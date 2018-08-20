package retry

import (
	"container/list"
	"time"
	"fmt"
)
type Retry struct{
	Max int
	Count int
	WaitingTime time.Duration
	Errors []*error
}

func (this *Retry) Error() string {
	str := ""
	for _, err := range this.Errors{
		str += fmt.Sprintf("%+v", err)
	}
	return str
}

func (this *Retry) Reset() {
	this.Count = 0
	this.Errors = []*error{}
}

type Retries struct{
	List *list.List // Retry
	Element *list.Element
	Retry *Retry
}

func (this *Retries) Error() string {
	str := ""
	for e := this.List.Front(); e != nil; e = e.Next() {
		str += e.Value.(*Retry).Error()
	}
	return str
}

func (this *Retries) Reset() {
	for e := this.List.Front(); e != nil; e = e.Next() {
		e.Value.(*Retry).Reset()
	}
	this.Element = this.List.Front()
	this.setRetry()
}
func (this *Retries) setRetry() bool {
	if this.Element != nil{
		this.Retry = this.Element.Value.(*Retry)
		return true
	}
	return false
}
func (this *Retries) Next() bool {
	this.Element = this.Element.Next()
	return this.setRetry()
}

func NewRetries(args ... interface{}) *Retries {
	l := len(args)
	retries := &Retries{
		List:list.New(),
	}
	for i := 0; i < l; i +=2 {
		r := &Retry{
			Max:args[i].(int),
			WaitingTime:args[i + 1].(time.Duration),
		}
		retries.List.PushBack(r)
	}
	retries.Reset()
	return retries
}