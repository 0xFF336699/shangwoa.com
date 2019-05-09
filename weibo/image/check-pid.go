package image

import (
	"net/http"
	"net/url"
)

func CheckPic(pid string) (err error, disabled bool) {
	u:="https://ww1.sinaimg.cn/large/" + pid
	res, err := http.Get(u)
	if err != nil{
		return
	}
	if res.StatusCode != 200{
		disabled = true
		return
	}
	if res.ContentLength == 8844 && res.Request.URL.Path == "/images/default_d_h_large.gif"{
		disabled = true
		return
	}

	u2, err := url.Parse(u)
	if err != nil{
		return
	}
	if u2.Path != res.Request.URL.Path{
		disabled = true
		return
	}
	return
}
