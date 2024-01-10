package rabbitmq

import (
	"github.com/streadway/amqp"
)

/**
PublishV2("amp.fanout", "express_sms_conf_change", "amqp://u1:pw@x.x.com:5674/virtual_host_name",
		publishing:= amqp.Publishing{
			Headers:         amqp.Table{"x-delay": 3000},
			ContentType:     "text/plain",
			DeliveryMode:    2,
			Body:            []byte("message body"),
		})
 */
func PublishV2(exchange, routingKey, url string, p amqp.Publishing) error {
	//publishing:= amqp.Publishing{
	//	Headers:         amqp.Table{"x-delay": 3000},
	//	ContentType:     "text/plain",
	//	ContentEncoding: "",
	//	DeliveryMode:    2,
	//	Priority:        0,
	//	CorrelationId:   "",
	//	ReplyTo:         "",
	//	Expiration:      "",
	//	MessageId:       "",
	//	Timestamp:       time.Time{},
	//	Type:            "",
	//	UserId:          "",
	//	AppId:           "",
	//	Body:            []byte("message body"),
	//}
	ch, err := GetChannelByDefaultRetries(url)
	if err != nil {
		return err
	}
	err = ch.Publish(exchange, routingKey, false, false, p)
	if err != nil {
		return err
	}
	return nil
}

func GetChan(url string, declare func(ch *amqp.Channel)error) error {
	ch, err := GetChannelByDefaultRetries(url)
	if err != nil {
		return err
	}
	return declare(ch)
}