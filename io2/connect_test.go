package io2

import (
	"fmt"
	"testing"
	"time"

	"strconv"

	"github.com/pkg/errors"
)

func TestConnector_GetClient(t *testing.T) {
	retryCount := 0
	connector := &Connector{
		Connect: func() (interface{}, error) {
			time.Sleep(time.Second * 1)
			retryCount++
			if retryCount == 3 || retryCount == 5 {
				return t, nil
			}
			return nil, errors.New("waiting" + strconv.Itoa(retryCount))
		},
		OnError: func(e error) bool {
			fmt.Println("OnError", retryCount)
			if retryCount == 3 {
				return true
			}
			//fmt.Println("error is", e)
			return false
		},
		RetryDelay: time.Second * 1,
	}
	forever := make(chan bool)
	for i := 0; i < 3; i++ {
		//fmt.Println("id", i)
		go func(index int) {
			//fmt.Println("index is", index)
			client, ch := connector.GetClient(false)
			if client != nil {
				fmt.Println("client is", client)
			} else {
				c := <-ch
				fmt.Println(c)
			}
		}(i)
	}
	go func() {
		time.Sleep(time.Second * 5)
		client, ch := connector.GetClient(false)
		if client != nil {
			fmt.Println("client is", client)
		} else {
			c := <-ch
			fmt.Println(c)
		}
	}()
	go func() {
		time.Sleep(time.Second * 7)
		client, ch := connector.GetClient(true)
		if client != nil {
			fmt.Println("client is", client)
		} else {
			c := <-ch
			fmt.Println(c)
		}
	}()
	go func() {
		time.Sleep(time.Second * 10000)
	}()
	<-forever
}
