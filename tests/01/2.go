package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"shangwoa.com/rabbitmq"
	"time"
)

func main()  {
	f:= make(chan int)
	go tListen()
	tPublishByConf()
	<-f
}

func tListen() {
	exec:= func(d *amqp.Delivery) error {
		fmt.Println("listen msg is", string(d.Body))
		return nil
	}
	onError:= func(qname string, err error, extra string) bool {
		fmt.Println("on error is", qname, err, extra)
		return false
	}
	url:= "amqp://album:album@rabbitmq.hb.ms.shangwoa.com:8231/album"
	qname := "crawl_url:insert"
	fmt.Println("qname is", qname)
	qd := &rabbitmq.QueueDeclare{Queue: qname, Durable: true, AutoDelete:false}
	qb := &rabbitmq.QueueBind{Queue: qname, RoutingKey: qname, Exchange: "amq.direct"}
	qos := &rabbitmq.Qos{PrefetchCount: 1, PrefetchSize: 0}
	consume := &rabbitmq.Consume{Queue: qname, ConsumerTag: ""}
	rabbitmq.Listen(exec, onError,url,qd, qb, qos, consume,time.Duration(60), 0, 20)
}

func tPublishByConf() {
	return
	body := []byte("hehe")
	url := "amqp://album:album@rabbitmq.hb.ms.shangwoa.com:8231/album"
	qname := "crawl_url:insert"
	routingKey := qname
	exchange := "amq.direct"
	kind := "direct"
	p := amqp.Publishing{DeliveryMode: 2, ContentType:  "text/plain", Body:body}
	declare := func(ch *amqp.Channel) (err error) {return ch.ExchangeDeclare(exchange , kind , true , false , false , false, nil)}
	declare = func(ch *amqp.Channel) (err error) {
		_, err = ch.QueueDeclare(qname, true, false, false, false, nil)
		return ch.ExchangeDeclare(exchange,   // name
			kind, // type direct, topic, headers and fanout
			true,     // durable
			false,    // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		)
	}
	publish := func(ch *amqp.Channel, p *amqp.Publishing)(err error) {return  ch.Publish(exchange, routingKey, false, false, *p)}
	err := rabbitmq.Publish2(url, &p, declare, publish)
	fmt.Println("publish error is", err)
}