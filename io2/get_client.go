package io2

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

// Connector
type Connector struct {
	Connect       func() (interface{}, error)
	OnError       func(error) bool
	OnMaxSeqAlarm func(int)
	MaxSeq        int
	RetryDelay    time.Duration
	waitingList   []chan<- interface{}
	client        interface{}
	isConnecting  bool
	lock          sync.RWMutex
	retryCount    int
}

// 如果有了客户端就返回客户端，没有就返回管道
// reconnect 是否注销已获取的client。因为有时候有些链接会失效，例如rabbitmq链接，所以这里加上强制重新链接
func (this *Connector) GetClient(reconnect bool) (interface{}, <-chan interface{}) {
	ch := make(chan interface{})
	if !reconnect && !this.isConnecting && this.client != nil {
		return this.client, nil
	} else {
		this.lock.Lock()
		this.client = nil
		this.waitingList = append(this.waitingList, ch)
		if this.OnMaxSeqAlarm != nil {
			l := len(this.waitingList)
			if l > this.MaxSeq {
				this.OnMaxSeqAlarm(l)
			}
		}

		this.lock.Unlock()
		go this.connect()
	}
	return this.client, ch
}
func (this *Connector) connect() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.isConnecting {
		return
	}
	if this.client != nil {
		go this.publish()
		return
	}
	this.retryCount++
	//fmt.Println("this is connecting"+strconv.Itoa(this.retryCount), this.client == nil)
	this.isConnecting = true
	client, err := this.Connect()
	//fmt.Println("error is", err)
	this.isConnecting = false
	if err != nil {
		stop := this.OnError(err)
		if stop {
			fmt.Println("break")
			return
		} else {
			go this.later()
			return
		}
	}
	this.client = client
	//fmt.Println("got client", this.retryCount)
	go this.publish()
}
func (this *Connector) publish() {
	this.lock.Lock()
	for _, w := range this.waitingList {
		w <- this.client
	}
	this.waitingList = nil
	this.lock.Unlock()
}

func (this *Connector) later() {
	time.Sleep(this.RetryDelay)
	go this.connect()
}
