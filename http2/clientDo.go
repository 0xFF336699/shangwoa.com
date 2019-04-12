package http2

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"shangwoa.com/image2"
	"time"
)
var UserAgent = `Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36`

type FileDownload func(url, p string)(err error, path string, w, h int)
func GetReq(url string) (r *http.Request, err error) {
	r, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("User-Agent", UserAgent)
	return
}
var proxy func(req *http.Request)(url *url.URL, err error)
func init()  {
	proxyUrl := os.Getenv("proxy_url")
	if len(proxyUrl) > 0{
		//proxyUrl, err := url.Parse("http://127.0.0.1:1080")
		proxyUrl, err := url.Parse(proxyUrl)
		if err != nil{
			panic(err)
		}
		proxy = http.ProxyURL(proxyUrl)
	}
}
func GetClient(useProxy bool) *http.Client {
	Jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: Jar}
	if useProxy {
		client.Transport =&http.Transport{Proxy: proxy}
	}
	return client
}

func GetBodyMultiTime(url string, useProxy bool, maxRepeatCount int, repeatDelaySecond int) (err error, res *http.Response) {
	repeatIndex := 0;
	for repeatIndex < maxRepeatCount{
		err, res = GetBody(url, useProxy)
		if err != nil && repeatIndex < maxRepeatCount{
			time.Sleep(time.Duration(repeatDelaySecond) * time.Second)
			repeatIndex ++
		}else{
			return
		}
	}
	return
}
func GetBody(url string, useProxy bool)(err error, res *http.Response)  {

	r, err := GetReq(url)
	if err != nil{
		return
	}
	client := GetClient(useProxy)
	res, err = client.Do(r)
	return
}
func GetJarClient() *http.Client {
	Jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: Jar}
	return client
}

func GetEmptyClient() *http.Client {
	return &http.Client{}
}

func ClientDo(r *http.Request, client *http.Client) (content []byte, err error) {
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	content, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return
}

func FileDownloader(filename, url string, useProxy bool) (err error, file string, w, h int) {
	//filename = "c:\\x.jpg"
	file = filename
	err, res := GetBody(url, useProxy)
	if err != nil{
		return
	}
	defer res.Body.Close()
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		return
	}
	_, err = io.Copy(f, res.Body)
	if err !=nil{
		return
	}
	w, h = image2.GetDimensions(f)
	return
}

func NewFileDownloader(useproxy bool)(FileDownload)  {
	return func(url, filename string)(err error, path string, w ,h int){
		return FileDownloader(filename, url, useproxy)
	}
}