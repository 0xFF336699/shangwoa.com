package time2

import "time"
const (
	FormatOnlyTime = "2006-01-02 15:04:05"
)
var CNLocation *time.Location
func init() {
	CNLocation, _= time.LoadLocation("Asia/Chongqing")
}