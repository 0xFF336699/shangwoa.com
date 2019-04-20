package image

import (
	"fmt"
	"os"
	"path/filepath"
	"shangwoa.com/http2"
	"time"
)

type GetCookies func(old string)(err error, cookies string)
type DownUploader struct {
	cookies string
	cookieIsAvailable bool
	tempPath string
	loadingCookieRetryCount int
	isLoadingCookies bool
	GetCookies GetCookies// 从shangwoa那边注入进来 file-down-upload.go
	DownLoader http2.FileDownload
	waitingUrls []*WaitingUrl
	uploader *ImageUploader
}
type WaitingUrl struct {
	filename string
	url      string
	ch       chan *WaitingUrl
	Err      error
	Pid      string
	Width int
	Height int
}
func NewDownUploader(tempPath string, getCookies GetCookies, downloader http2.FileDownload) (loader *DownUploader) {
	loader = &DownUploader{
		tempPath:tempPath,
		GetCookies:getCookies,
		DownLoader:downloader,
		waitingUrls: []*WaitingUrl{},
		uploader:NewUploader(""),
	}
	return
}
func (this *DownUploader)loadCookies()  {
	if this.isLoadingCookies {
		return
	}
	this.loadingCookieRetryCount = 0
	this.isLoadingCookies = true
	this.getCookies()
	this.isLoadingCookies = false

}
func (this *DownUploader) getCookies() (err error) {
	this.cookieIsAvailable = false
	err, cookies := this.GetCookies(this.cookies)
	if err != nil && this.loadingCookieRetryCount < 5{
		this.loadingCookieRetryCount ++
		time.Sleep(time.Duration(60) * time.Second)
		return this.getCookies()
	}
	if err == nil && len(cookies) > 100{
		this.cookies = cookies
		this.cookieIsAvailable = true
		this.uploader.cookie = cookies
		go this.checkWorkOk()
		return
	}
	if err != nil{
		this.cookies = ""
		return
	}
	return nil
}
func (this *DownUploader)DownUpload(url, filename string) (c <-chan *WaitingUrl) {
	ch := make(chan *WaitingUrl)

	p := filepath.Join(this.tempPath , "/" + filename)
	this.waitingUrls = append(this.waitingUrls, &WaitingUrl{url: url,filename:p,ch:ch})
	go this.checkWorkOk()
	return ch
}
func (this *DownUploader) checkWorkOk() {
	if this.isLoadingCookies == false && len(this.cookies) > 0{
		for _, w :=range this.waitingUrls{
			go this.downUpload(w)
		}
	}else{
		this.loadCookies()
	}
}

func (this *DownUploader)downUpload(w *WaitingUrl)  {
	defer this.waitingUrlWorkDown(w)
	fmt.Println("downupload start ", w.url)
	err, u:= this.DownLoader(w.url, w.filename)
	if err != nil{
		fmt.Println("weibo downloader error", w.url, err.Error())
		w.Err = err
		w.Pid = u
		w.ch <- w
		return
	}

	err, pid, width, height := this.uploader.Upload(u)
	if err != nil{
		fmt.Println("weibo upload error", w.url, err.Error())
		this.cookieIsAvailable = false
		go this.loadCookies()
		w.Err = err
		w.Pid = pid
		w.ch <- w
		return
	}
	w.Pid = pid
	w.Width = width
	w.Height = height
	fmt.Println("downupload complete", w.url)
	w.ch <- w
	return
}

func (this *DownUploader) waitingUrlWorkDown(w *WaitingUrl)  {
	go removeFile(w.filename)
	index := -1
	for i, v := range this.waitingUrls{
		if v == w{
			index = i
			break;
		}
	}
	if index > -1{
		this.waitingUrls = append(this.waitingUrls[:index], this.waitingUrls[index + 1:]...)
	}
}

func removeFile(f string) (err error) {
	time.Sleep(1000)
	err = os.Remove(f)
	fmt.Println("down-upload removeFile", f, err)
	return
}