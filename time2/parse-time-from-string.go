package time2

import (

	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)
func ParseTimeFromString(str string)(err error, isOk bool, t time.Time)  {
	if str == `""`{
		return
	}
	if string(str[0]) == `"`{
		str = str[1:]
	}
	if string(str[len(str) - 1]) == `"`{
		str = str[:len(str) - 1]
	}
	if str == "null" {
		return
	}
	if len(str) == 0{
		return
	}
	if ok, err := regexp.Match("^(\\d+)$", []byte(str)); err == nil && ok{
		millis, err := strconv.ParseInt(str, 10, 64)
		if err != nil{
			return err, false, t
		}
		t = time.Unix(0, millis * int64(time.Millisecond))
		isOk = true
		return err, isOk, t
	}
	layout := ""
	if len(str) == 10{
		layout = "2006-01-02"
	}else if len(str) >= 19{
		if index := strings.Index(str, "T"); index == -1{
			layout = "2006-01-02 15:04:05"
		}else{
			layout = "2006-01-02T15:04:05"
		}
		if len(str) > 19{
			remainder := str[19:]
			reg, _ := regexp.Compile(`\d`)
			remainder = reg.ReplaceAllString(remainder, "0")
			layout += remainder
		}
	}
	if len(layout) == 0{
		err = errors.New("unknow input type")
		return
	}
	t, err = time.Parse(layout, str)
	if err != nil{
		return
	}
	isOk = true
	return
}

