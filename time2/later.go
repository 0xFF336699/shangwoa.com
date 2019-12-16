package time2

import (
	"github.com/pkg/errors"
	"time"
)

var LaterErrorInvalidAfter = errors.New("invalid after")
var LaterErrorNoCallback = errors.New("no callback");
type callback func()
type Later struct{
	after time.Duration
	t *time.Timer
	cb callback
	c chan bool
	canceled bool
	completed bool
}

func (this *Later) Start()  {
	defer this.t.Stop()
	for {
		select {
		case <-this.t.C:
			this.cb()
			this.completed = true
			goto END
			case <-this.c:
				goto END
		}
	}
	END:
}
func (this *Later) Cancel()(bool)  {
	if this.canceled{
		return false
	}
	if this.completed{
		return false
	}
	this.c <- true
	this.canceled = true
	return this.t.Stop()
}
func NewLater(after time.Duration, cb callback, start bool) (err error, later *Later) {
	if after <= 0{
		return LaterErrorInvalidAfter, nil
	}
	t := time.NewTimer(after)
	if cb == nil{
		return LaterErrorNoCallback, nil
	}
	later = &Later{after,t, cb, make(chan bool, 1), false, false}
	if start{
		go later.Start()
	}
	return
}