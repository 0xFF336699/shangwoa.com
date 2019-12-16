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
	// https://ww1.sinaimg.cn/large/007aNjkZly1g24kbnr405j30fa0mwmyi
	url = "https://ww" + strconv.Itoa(i) + ".sinaimg.cn/" + size + "/" + pid
	return
}

func GetImageUrl(pid, size string) ( url string) {
	i, err := number.Rand(1, 4)
	if err != nil{
		i = 1
	}
	// https://ww1.sinaimg.cn/large/007aNjkZly1g24kbnr405j30fa0mwmyi
	url = "https://ww" + strconv.Itoa(i) + ".sinaimg.cn/" + size + "/" + pid
	return
}