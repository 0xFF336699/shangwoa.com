package main

import (
	"fmt"
	"shangwoa.com/rabbitmq"
	"time"
)

func main() {
	f := make(chan int)
	go createConsumePublish()
	<-f
}

func createConsumePublish() {
	testChannel()
}

func testChannel() {
	const(
		//testChannelExchange = "amq.fanout"
		testChannelKind = "fanout"
		//testChannelExchange = "amq.direct"
		testChannelExchange = "order-paid"
		//testChannelKind = "direct"
		testChannelQname = "order-paid"
		//URL = "amqp://album:album@rabbitmq.hb.ms.shangwoa.com:8231/album"
		URL = "amqp://snacks:oNQHPaSuCsblb6MX@rabbitmq.hb.ms.fanfanlo.com:8231/snacks-v"
	)
	args := &rabbitmq.ChannelArgs{
		URL:URL,
		Exchange:testChannelExchange,
		Qname:testChannelQname,
		Kind:testChannelKind,
		RoutingKey:testChannelQname,
		Durable:true,
		Mandatory:false,
		Immediate:false,
		AutoDelete:false,
		Internal:false,
		Exclusive:false,
		NoWait:false,
		RetryMaxCount:10,
		InteralTime:1000,
		Args:nil,
		DeliveryMode:2,
		ContentType:"text/plain",
		MaxWaitingCount:10,
	}

	err, pc := rabbitmq.MaxWaitingPublishChannel(args)
	if err!= nil{
		panic(err)
	}
	testChannelPublish(pc, []byte("hehe"))
	testChannelPublish(pc, []byte("abc"))
	count := 0
	for{
		err = testChannelPublish(pc, []byte(time.Now().Format(time.RFC3339Nano)))
		if err != nil{
			time.Sleep(time.Second * 5)
			go testChannel()
			return
		}
		count ++
		w := time.Millisecond * 200
		if count%20 == 1{
			w = time.Second * 5
		}
		time.Sleep(w)
	}
}

func testChannelPublish(pc *rabbitmq.PubChannel, body []byte)(err error){
	err = pc.Publish(body)
	if err != nil{
		fmt.Println("err is ", err.Error())
	}
	return
}