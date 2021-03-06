package rabbitmq

import (
	"github.com/streadway/amqp"
	"testing"
	"fmt"
	"time"
)
//go test -v publish_test.go publish.go rabbitmq.go
func TestPublishByDefault(t *testing.T) {
	//return
	body := []byte("hehe")
	err := PublishByDefault("test:1", "amqp://album:album@rabbitmq.hb.ms.shangwoa.com:8231/album", body)
	if err != nil{
		fmt.Println(err)
	}
}

func TestListen(t *testing.T) {
	exec:= func(d *amqp.Delivery) error {
		fmt.Println("listen msg is", string(d.Body))
		return nil
	}
	onError:= func(qname string, err error, extra string) bool {
		fmt.Println("on error is", qname, err, extra)
		return false
	}
	url:= "amqp://album:album@rabbitmq.hb.ms.shangwoa.com:8231/album"
	qname := "test:2"
	fmt.Println("qname is", qname)
	qd := &QueueDeclare{Queue: qname, Durable: true, AutoDelete:false}
	qb := &QueueBind{Queue: qname, RoutingKey: qname, Exchange: "amq.direct"}
	qos := &Qos{PrefetchCount: 1, PrefetchSize: 0}
	consume := &Consume{Queue: qname, ConsumerTag: ""}
	Listen(exec, onError,url,qd, qb, qos, consume,time.Duration(60), 0, 20)
}

func TestPublishByConf(t *testing.T) {
	body := []byte("hehe")
	url := "amqp://album:album@rabbitmq.hb.ms.shangwoa.com:8231/album"
	qname := "test:2"
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
	Publish2(url, &p, declare, publish)
}