package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"shangwoa.com/rabbitmq"
	"time"
)

func main()  {
	f := make(chan int)
	go consume()
	//testChannel()
	tickOldPublish()
	<-f
}

func consume() {

	qname := "test"
	qd := &rabbitmq.QueueDeclare{Queue: qname, Durable: true,AutoDelete:false}
	qb := &rabbitmq.QueueBind{Queue: qname, RoutingKey: qname, Exchange: "amq.direct"}
	qos := &rabbitmq.Qos{PrefetchCount: 1, PrefetchSize: 0}
	consume := &rabbitmq.Consume{Queue: qname, ConsumerTag: ""}
	fmt.Println("listen post ", qname, "amqp://album:album@rabbitmq.hb.ms.shangwoa.com:8231/album")
	rabbitmq.Listen(onMessage, onError, "amqp://album:album@rabbitmq.hb.ms.shangwoa.com:8231/album", qd, qb, qos, consume, time.Duration(60),0, 20)
}
func onMessage(d *amqp.Delivery) (err error) {
	fmt.Println("on message", string(d.Body))
	return
}

func onError(qname string, err error, extra string) bool {
	fmt.Println(qname, err, extra)
	if err.Error() == rabbitmq.EmptyBodyError.Error() {
		return true
	}
	switch err.Error() {
	case rabbitmq.EmptyBodyError.Error():
		return true
	case rabbitmq.RetryAlarmError.Error():
		return true
	}
	return false
}
func testChannel() {
	args := &rabbitmq.ChannelArgs{
		URL:"amqp://album:album@rabbitmq.hb.ms.shangwoa.com:8231/album",
		Exchange:"amq.direct",
		Qname:"test",
		Kind:"direct",
		RoutingKey:"test",
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
}

func testChannelPublish(pc *rabbitmq.PubChannel, body []byte){
	err := pc.Publish(body)
	if err != nil{
		fmt.Println("err is ", err.Error())
	}
}

func tickOldPublish() {
	t := time.NewTicker(time.Second * 1)
	for{
		<-t.C
		oldPublish()
	}
}
func oldPublish() {

	qname := "test"
	err := rabbitmq.Publish(qname, "amq.direct", qname, "direct", "amqp://album:album@rabbitmq.hb.ms.shangwoa.com:8231/album", []byte(time.Now().String()))
	if err != nil{
		fmt.Println("old publish err is", err.Error())
	}else{
		fmt.Println("publishedold publish")
	}
}