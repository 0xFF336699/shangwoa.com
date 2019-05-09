package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"shangwoa.com/mq/rabbit"
	"time"
)

func main() {
	f := make(chan bool)
	go createConsumer()
	//createConsumePublish()
	<- f
}
const(
	//testChannelExchange = "amq.fanout"
	//testChannelKind = "fanout"
	testChannelExchange = "amq.direct"
	testChannelKind = "direct"
	testChannelQname = "test"
	URL = "amqp://album:album@rabbitmq.hb.ms.shangwoa.com:8231/album"
)

func createConsumer() {
	args := &rabbit.ConsumerInitChannelArgs{Exchange:"amq.direct",Qname:"test",RoutingKey:"test", Durable:true, AutoDelete:false, Exclusive:false, NoWait:false, PrefetchCount:1, PrefetchSize:0, Global:false, AutoAck:false, NoLocal:false}
	connection := &rabbit.Connection{Url:URL,RetryMaxCount:10}
	channel := &rabbit.Channel{RetryMaxCount:10, Connection:connection}
	consumer:= &rabbit.Consumer{Channel:channel, OnGetChannel:rabbit.CreateOnGetChannelByArgs(args), OnGetMessage:consumerOnGetMessage}
	consumer.OnLaunchError =  func(err error) {
		time.Sleep(time.Second * 10)
		consumer.Start()
	}
	consumer.Start()
}
var count = 0
var emptyCount = 0
func consumerOnGetMessage(d *amqp.Delivery) (handled bool, closeChannel, closeConnection bool) {
	fmt.Println("body is", count, string(d.Body))
	handled = true
	count ++
	if len(d.Body) == 0{
		emptyCount ++
		fmt.Println("---------------------", emptyCount)
		if emptyCount > 3{
			closeConnection = true
			emptyCount = 0
			fmt.Println("restart xxx")
			return
		}
	}
	if count == 5{
		closeChannel = true
	}
	if count == 10{
		closeConnection = true
		count = 0
	}
	return
}
func consumerOnGetChannel(ch *amqp.Channel) (err error, msgs <-chan amqp.Delivery) {
	count := 0
	for{
		_, err = ch.QueueDeclare(testChannelQname,true, false, false, false, nil)
		if err == nil{
			break
		}
		count ++
		if count > 10{
			return
		}
		time.Sleep(time.Second * 1 * time.Duration(count))
	}

	count = 0
	for {
		err = ch.QueueBind(testChannelQname, testChannelQname, testChannelExchange, false, nil)
		if err == nil{
			break
		}
		count ++
		if count > 10{
			return
		}
		time.Sleep(time.Second * 1 * time.Duration(count))
	}

	count = 0
	for {
		err = ch.Qos(1, 0, false)
		if err == nil{
			break
		}
		count ++
		if count > 10{
			return
		}
		time.Sleep(time.Second * 1 * time.Duration(count))
	}
	count = 0
	for {
		msgs, err = ch.Consume(testChannelQname, "", false, false, false, false, nil)
		if err == nil{
			break
		}
		count ++
		if count > 10{
			return
		}
		time.Sleep(time.Second * 1 * time.Duration(count))
	}
	return
}