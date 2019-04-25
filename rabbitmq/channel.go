package rabbitmq

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"shangwoa.com/utils/retry"
	"sync/atomic"
	"time"
)
type ChannelArgs struct {
	URL string
	Exchange string
	Qname string
	Kind string
	RoutingKey string
	Durable bool
	Mandatory bool
	Immediate bool
	AutoDelete bool
	Internal bool
	Exclusive bool
	NoWait bool
	RetryMaxCount int
	InteralTime int
	Args amqp.Table
	DeliveryMode uint8
	ContentType string
	MaxWaitingCount int32
}
var WaitingFlowError = errors.New("waiting flow")
type MaxWaitingChannel struct {
	ch              *amqp.Channel
	waitingCount    int32
	args *ChannelArgs
}
type PubChannel struct {
	MaxWaitingChannel
}
func (this *PubChannel)Publish(body []byte)(err error){
	args := this.args
	c := atomic.LoadInt32(&this.waitingCount)
	if c >= args.MaxWaitingCount{
		return WaitingFlowError
	}
	n := c + 1
	for{
		if atomic.CompareAndSwapInt32(&this.waitingCount, c, n){
			break
		}
	}
	fmt.Println("rabbitmq publish befor atomic is", n)
	isClosed := false
	err, _, _ = retry.RetryDoInteralFunc(func() (err error, res interface{}) {
		err = this.ch.Publish(args.Exchange, args.RoutingKey, args.Mandatory, args.Immediate, amqp.Publishing{DeliveryMode:args.DeliveryMode,ContentType:args.ContentType, Body:body})
		if err == amqp.ErrClosed{
			isClosed = true
		}
		return
	}, func(count int, err error) (jump bool) {
		if count >= args.RetryMaxCount || isClosed{
			return true
		}
		time.Sleep(time.Duration(count * args.InteralTime))
		return false
	})
	n = atomic.AddInt32(&this.waitingCount, -1)
	fmt.Println("rabbitmq after publish atomic is", n)
	return
}

func PublishChannel(args *ChannelArgs)(err error, ch *amqp.Channel)  {
	err, res, _ := retry.RetryDoInteralTimeIncr(func() (err error, res interface{}) {
		res, err = GetChannelByDefaultRetries(args.URL)
		return
	}, args.RetryMaxCount, args.InteralTime)
	if err != nil{
		fmt.Println("PublishChannel error is", err.Error())
		return
	}
	ch = res.(*amqp.Channel)
	if ch == nil{
		fmt.Println("PublishChannel not a channel")
		return errors.New("not a channel"), nil
	}
	err, _, _ = retry.RetryDoInteralTimeIncr(func() (err error, res interface{}) {
		return ch.ExchangeDeclare(args.Exchange , args.Kind , args.Durable , args.AutoDelete , args.Internal , args.NoWait, args.Args), nil
	}, args.RetryMaxCount, args.InteralTime)
	if err != nil{
		return
	}

	err, _, _ = retry.RetryDoInteralTimeIncr(func() (err error, res interface{}) {
		res, err = ch.QueueDeclare(args.Qname, args.Durable, args.AutoDelete, args.Exclusive, args.NoWait, args.Args)
		return
	}, args.RetryMaxCount, args.InteralTime)
	if err != nil{
		return
	}
	//err, _, _ = retry.RetryDoInteralTimeIncr(func() (err error, res interface{}) {
	//	return ch.QueueBind(args.Qname, args.RoutingKey, args.Exchange, args.NoWait, args.Args), nil
	//}, args.RetryMaxCount, args.InteralTime)
	return
}

func MaxWaitingPublishChannel(args *ChannelArgs) (err error, mc *PubChannel) {
	err, ch:= PublishChannel(args)
	if err != nil{
		return
	}
	mc = &PubChannel{MaxWaitingChannel{ch:ch, args:args}}
	atomic.StoreInt32(&mc.waitingCount, 0)
	return
}