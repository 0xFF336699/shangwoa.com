package rabbitmq

import (
	"time"

	"github.com/streadway/amqp"
	"shangwoa.com/system"
)

func PublishByDefault(qName, url string, body []byte) error {
	//url := "amqp://ig-crawler:ig-crawler@rabbitmq.hb.ms.shangwoa.com:8231/ig-crawler"
	ch, err := GetChannelByDefaultRetries(url)
	if err != nil {
		return err
	}

	//qName := "post_media_order:insert"
	_, err = ch.QueueDeclare(qName, true, false, false, false, nil)
	if err != nil {
		retries := NewRetries(3, 1*time.Second, 2, 2*time.Second)
		for retries.Retry != nil && retries.Retry.Count < retries.Retry.Max {
			retries.Retry.Errors = append(retries.Retry.Errors, &err)
			retries.Retry.Count++
			time.Sleep(retries.Retry.WaitingTime)
			_, err = ch.QueueDeclare(qName, true, false, false, false, nil)
			if err != nil {
				if retries.Retry.Count >= retries.Retry.Max {
					bl := retries.Next()
					if !bl {
						e := &system.Error{Type: "DefaultPublish.Publish", Err: err, SubType: "", SubData: retries, CodeLevel: 1}
						return e
					}
				}
			}
		}
	}
	p := amqp.Publishing{
		DeliveryMode: 2,
		ContentType:  "text/plain",
		Body:         body,
	}

	err = ch.Publish(qName, qName, false, false, p)
	if err != nil {
		retries := NewRetries(3, 1*time.Second, 2, 2*time.Second)
		for retries.Retry != nil && retries.Retry.Count < retries.Retry.Max {
			retries.Retry.Errors = append(retries.Retry.Errors, &err)
			retries.Retry.Count++
			time.Sleep(retries.Retry.WaitingTime)
			err = ch.Publish(qName, qName, false, false, p)
			if err != nil {
				if retries.Retry.Count >= retries.Retry.Max {
					bl := retries.Next()
					if !bl {
						e := &system.Error{Type: "DefaultPublish.Publish", Err: err, SubType: "ch.Publish(", SubData: retries, CodeLevel: 1}
						return e
					}
				}
			}
		}
	}
	return nil
}


func Publish(qName, exchange, key, kind, url string, body []byte) error {
	//url := "amqp://ig-crawler:ig-crawler@rabbitmq.hb.ms.shangwoa.com:8231/ig-crawler"
	ch, err := GetChannelByDefaultRetries(url)
	if err != nil {
		return err
	}

	//qName := "post_media_order:insert"
	//_, err = ch.QueueDeclare(qName, true, false, false, false, nil)
	//if err != nil {
	//	return err
	//}
	p := amqp.Publishing{
		DeliveryMode: 2,
		ContentType:  "text/plain",
		Body:         body,
	}
	err = ch.ExchangeDeclare(
		exchange,   // name
		kind, // type direct, topic, headers and fanout
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}

	//err = ch.QueueBind(
	//	qName,      // queue name
	//	key, // routing key
	//	exchange,   // exchange
	//	false,
	//	nil)

	err = ch.Publish(exchange, key, false, false, p)
	if err != nil {
		return err
	}
	return nil
}

func Publish2(url string, p *amqp.Publishing, exchangeDeclare func(ch *amqp.Channel)(error), publish func(ch *amqp.Channel, p *amqp.Publishing)(error))(err error)  {
	ch, err := GetChannelByDefaultRetries(url)
	if err != nil {
		return
	}
	err = exchangeDeclare(ch)
	if err != nil{
		return
	}
	err = publish(ch, p)
	if err != nil{
		return
	}
	return
}


func Publish3(ch *amqp.Channel, p *amqp.Publishing, exchangeDeclare func(ch *amqp.Channel)(error), publish func(ch *amqp.Channel, p *amqp.Publishing)(error))(err error)  {
	err = exchangeDeclare(ch)
	if err != nil{
		return
	}
	err = publish(ch, p)
	if err != nil{
		return
	}
	return
}
