package time2

import "time"
const (
	FormatOnlyTime = "2006-01-02 15:04:05"
	CNLocationString = "Asia/Shanghai"
	//CNLocationString = "Asia/Chongqing"
)
var CNLocation *time.Location
func init() {
	CNLocation, _= time.LoadLocation("Asia/Chongqing")
}