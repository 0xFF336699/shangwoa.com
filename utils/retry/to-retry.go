package retry

import (
	"errors"
	"time"
)

type exec func ()(error, interface{});
type onError func(err error, count int)bool

func CreateSleepTimeOnError(d time.Duration, maxCount int) onError {
	return func(err error, count int) bool {
		if count > maxCount {
			return false
		}
		time.Sleep(d)
		return true
	}
}
func ToRetry(exec exec, onError onError)(err error, res interface{})  {
	i := 0
	for ;;i++ {
		err, res = exec()
		if err == nil{
			return
		}
		c := onError(err, i)
		if !c {
			return
		}
	}
}

func testToRetry() {
	err, res := ToRetry(func() (err error, i interface{}) {
		return errors.New("test"), i
	}, CreateSleepTimeOnError(time.Second * 1, 20));
	print(err, res)
}