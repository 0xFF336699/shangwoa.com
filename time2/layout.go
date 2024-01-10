package time2

import (
	"fmt"
	"time"
)
const (
	FormatOnlyTime = "2006-01-02 15:04:05"
	FormatOnlyDate= "2006-01-02"
	FormatDateHour= "2006-01-02 15"
	FormatTzTime = "2006-01-02T15:04:05Z"
	CNLocationString = "Asia/Shanghai"
	//CNLocationString = "Asia/Chongqing"
)
var CNLocation *time.Location
func init() {
	var err error
	CNLocation, err = time.LoadLocation("Asia/Chongqing")
	if err != nil{
		fmt.Println("time2 init cnlocation err", err)
	}
}
