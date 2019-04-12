package rabbitmq

import (
	"github.com/streadway/amqp"
	"testing"
	"fmt"
)

func TestPublishByDefault(t *testing.T) {
	body := []byte("hehe")
	err := PublishByDefault("post_media_order:downloaded", "amqp://ig-crawler:ig-crawler@rabbitmq.hb.ms.shangwoa.com:8231/ig-crawler", body)
	if err != nil{
		fmt.Println(err)
	}
}

func TestPublishByConf(t *testing.T) {
	body := []byte("hehe")
	url := "amqp://ig-crawler:ig-crawler@rabbitmq.hb.ms.shangwoa.com:8231/ig-crawler"
	qname := "xx"
	routingKey := qname
	exchange := "amq.direct"
	kind := "direct"
	p := amqp.Publishing{DeliveryMode: 2, ContentType:  "text/plain", Body:body}
	declare := func(ch *amqp.Channel) (err error) {return ch.ExchangeDeclare(exchange , kind , true , false , false , false, nil)}
	declare = func(ch *amqp.Channel) (err error) {
		return ch.ExchangeDeclare(exchange , kind , true , false , false , false, nil)
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