package image

import (
	"shangwoa.com/utils/number"
	"strconv"
)
const (
	WeiboSizeLarge = "large"
)
func GetRandomUrl(pid, size string) (err error, url string) {
	i, err := number.Rand(1, 4)
	if err != nil{
		return
	}
	url = "https://ww" + strconv.Itoa(i) + ".sinaimg.cn/" + size + "/" + pid
	return
}
