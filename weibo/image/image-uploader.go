package image

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/pkg/errors"
)

type ImageUploader struct {
	cookie string
}

type JTSData struct {
	Code string `json:"code"`
	Data Data   `json:"data"`
}

type Data struct {
	Data  string `json:"data"`
	Pics  Pics   `json:"pics"`
	Count int    `json:"count"`
}

type Pics struct {
	Pic1 Pic1 `json:"pic_1"`
	Pic2 Pic2 `json:"pic_2"`
}

type Pic1 struct {
	Ret    int    `json:"ret"`
	Height int    `json:"height"`
	Name   string `json:"name"`
	Pid    string `json:"pid"`
	Width  int    `json:"width"`
	Size   int    `json:"size"`
}

type Pic2 struct {
	Ret  int    `json:"ret"`
	Name string `json:"name"`
}

func (this *ImageUploader) Upload(filePath string) (err error, pid string, w, h int) {

	//uploadURL := "http://picupload.service.weibo.com/interface/pic_upload.php?ori=1&mime=image%2Fjpeg&data=base64&url=0&markpos=1&logo=&nick=0&marks=1&app=miniblog"
	uploadURL := "http://picupload.service.weibo.com/interface/pic_upload.php?mime=image%2Fjpeg&data=base64&url=0&markpos=1&logo=&nick=0&marks=1&app=miniblog"
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	formFile, err := writer.CreateFormFile("pic1", filePath)
	if err != nil {
		log.Fatalf("Create form file failed: %s\n", err)
	}

	// 从文件读取数据，写入表单
	srcFile, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("%Open source file failed: s\n", err)
		return
	}
	c, _, err := image.DecodeConfig(srcFile)
	if err == nil{
		w = c.Width
		h = c.Height
	}
	srcFile.Close()
	time.Sleep(1)
	srcFile, err = os.Open(filePath)
	if err != nil {
		log.Fatalf("%Open source file failed: s\n", err)
		return
	}
	defer srcFile.Close()
	_, err = io.Copy(formFile, srcFile)
	if err != nil {
		log.Fatalf("Write to form file falied: %s\n", err)
		return
	}

	// 发送表单
	//contentType := writer.FormDataContentType()
	writer.Close() // 发送之前必须调用Close()以写入结尾行
	//_, err = http.Post("http://picupload.service.weibo.com/interface/pic_upload.php?ori=1&mime=image%2Fjpeg&data=base64&url=0&markpos=1&logo=&nick=0&marks=1&app=miniblog", contentType, buf)
	//if err != nil {
	//	log.Fatalf("Post failed: %s\n", err)
	//}
	req, err := http.NewRequest("POST", uploadURL, buf)
	if err != nil {
		fmt.Println("image-uuploader Upload err", err.Error())
		return
	}
	t := writer.FormDataContentType()
	//t = "blob"
	req.Header.Set("Content-Type", t)
	//req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Cookie", this.cookie)
	//userAgent := `Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36`
	//req.Header.Set("User-Agent", userAgent)
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		//fmt.Println("g")
		return
	}
	content, err := ioutil.ReadAll(res.Body)
	str := string(content)
	//io.Copy(os.Stderr, res.Body) // Replace this with Status.Code check
	reg, _ := regexp.Compile(`[\s\S]*script>\n`)
	s := reg.ReplaceAllString(str, "")
	data := JTSData{}
	err = json.Unmarshal([]byte(s), &data)
	if err != nil {
		return
	}
	pid = data.Data.Pics.Pic1.Pid
	if pid == "" {
		fmt.Println("upload return str is", str)
		err = errors.New("upload lost result is " + str)
	}
	//http://ww1.sinaimg.cn/large/
	return
}
func NewUploader(cookie string) (uploader *ImageUploader) {
	uploader = &ImageUploader{cookie: cookie}
	return uploader
}
