package broadcast

import (
	"reflect"
	"sync"
)

type callbackFunc func(i interface{})
type Broadcast struct{
	m sync.RWMutex
	chListeners []chan<- interface{}
	callbacks [] callbackFunc
}

func NewBroadcaster() *Broadcast {
	return &Broadcast{
		m:           sync.RWMutex{},
		chListeners: [] chan<- interface{}{},
		callbacks: []callbackFunc{},
	}
}

func (this *Broadcast) AddListener(ch chan<- interface{}) (inserted bool) {
	this.m.RLock()
	defer this.m.RUnlock()
	for _, v := range this.chListeners{
		if v == ch{
			return false
		}
	}
	this.chListeners = append(this.chListeners, ch)
	return true
}

func (this *Broadcast) RemoveListener(ch chan interface{}) (removed bool) {
	this.m.Lock()
	defer this.m.Unlock()
	for i, v := range this.chListeners{
		if v == ch{
			this.chListeners = append(this.chListeners[:i], this.chListeners[i + 1:]...)
			return true
		}
	}
	return false
}

func (this *Broadcast) AddCallback(f callbackFunc) (inserted bool) {
	this.m.RLock()
	defer this.m.RUnlock()
	v1 := reflect.ValueOf(f)
	for _, v := range this.callbacks{
		v2 := reflect.ValueOf(v)
		if v1 == v2{
			return false
		}
	}
	this.callbacks = append(this.callbacks, f)
	return true
}

func (this *Broadcast) RemoveCallback(f callbackFunc) (removed bool) {
	this.m.Lock()
	defer this.m.Unlock()
	v1 := reflect.ValueOf(f)
	for i, v := range this.chListeners{
		v2 := reflect.ValueOf(v)
		if v1 == v2{
			this.callbacks = append(this.callbacks[:i], this.callbacks[i + 1:]...)
			return true
		}
	}
	return false
}
// 侦听时需要根据情况来决定是否要对i进行加锁
func (this *Broadcast) Broadcast(i interface{}) {
	this.m.RLock()
	defer this.m.RUnlock()
	for _, v := range this.chListeners{
		v <- i
	}
	for _, v := range this.callbacks{
		go v(i)
	}
}