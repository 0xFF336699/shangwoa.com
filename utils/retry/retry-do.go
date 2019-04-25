package retry

import (
	"fmt"
	"time"
)

func RetryDoInteralTime(f func()(err error, res interface{}), retryMaxCount int, interalTime int) (err error, res interface{}, count int) {
	count = 0
	for{
		var e error
		e, res = f()
		if e != nil{
			fmt.Println("er r is", e.Error())
			count ++
			if count > retryMaxCount{
				err = e
				return
			}
			time.Sleep(time.Duration(interalTime))
		}
		if e == nil{
			return
		}
	}
}

func RetryDoInteralTimeIncr(f func()(err error, res interface{}), retryMaxCount int, interalTime int)  (err error, res interface{}, count int) {
	count = 0
	for{
		var e error
		e, res = f()
		if e != nil{
			count ++
			if count > retryMaxCount{
				err = e
				return
			}
			time.Sleep(time.Duration(interalTime * count))
		}
		if e == nil{
			return
		}
	}
}

func RetryDoInteralFunc(f func()(err error, res interface{}), retryFunc func(count int, err error)(jump bool)) (err error, res interface{}, count int) {
	count = 0
	for{
		var e error
		e, res = f()
		count ++
		if e != nil{
			if retryFunc(count, e){
				return
			}
		}
		if e == nil{
			return
		}
	}
}