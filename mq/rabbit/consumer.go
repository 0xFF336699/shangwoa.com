package rabbit

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

type Consumer struct{
	Channel       *Channel
	OnGetChannel  func(ch *amqp.Channel) (err error, msgs <-chan amqp.Delivery)
	OnGetMessage  func(d *amqp.Delivery)(handled bool, closeChannel, closeConnection bool)
	OnLaunchError func(err error)
	msgs          <-chan amqp.Delivery
	chClosed      [] chan bool
}
func (this *Consumer)GetChannel() (err error, ch *amqp.Channel) {
	err, ch = this.Channel.Channel()
	if err != nil{
		return
	}
	this.listenChannelClose()
	return
}

func (this *Consumer) listenChannelClose()  {
	ch := this.Channel.Ch
	receiver := make(chan *amqp.Error)
	receiver = ch.NotifyClose(receiver)
	go func() {
		for{
			select{
			case  e :=<-receiver:
				fmt.Println("listenChannelClose e is", e)
				this.Restart()
				return
			}
		}
	}()
	return
}
func (this *Consumer) Restart() {
	fmt.Println("consumer restart")
	for _, ch := range this.chClosed{
		ch <- true
	}
	this.chClosed = nil
	time.Sleep(time.Second * 1)
	this.Start()
}
func (this *Consumer)Start() {
	this.chClosed = []chan bool{}
	err, _ := this.GetChannel()
	if err != nil{
		this.OnLaunchError(err)
		return
	}
	err, this.msgs = this.OnGetChannel(this.Channel.Ch)
	if err != nil{
		this.OnLaunchError(err)
		return
	}
	this.Consume()
}
func (this *Consumer)Consume()  {
	emptyCount := 0
	fmt.Println("rabbit mq consumer working")
	for{
		select{
		case d, ok := <-this.msgs:
			if !ok{
				time.Sleep(time.Second * 2)
				this.Restart()
			}
			if len(d.Body) == 0{
				emptyCount ++
				d.Ack(false)
				if emptyCount > 5{
					go this.Restart()
					return
				}
				break
			}else{
					emptyCount = 0
			}
			handled, closeChannel, closeConnection:= this.OnGetMessage(&d)
			if handled{
				d.Ack(false)
			}
			if closeChannel{
				time.Sleep(time.Second * 2)
				go this.Restart()
				return
			}
			if closeConnection{
				time.Sleep(time.Second * 2)
				go this.Channel.Connection.Restart()
				return
			}

		}
	}
	fmt.Println("rabbit mq consumer workend")
}


type ConsumerInitChannelArgs struct {
	Exchange string
	Qname string
	RoutingKey string
	Durable bool
	AutoDelete bool
	Exclusive bool
	NoWait bool
	Args amqp.Table
	PrefetchCount int
	PrefetchSize int
	Global bool
	Consumer string
	AutoAck bool
	NoLocal bool
}
func CreateOnGetChannelByArgs(args *ConsumerInitChannelArgs) (func (ch *amqp.Channel) (err error, msgs <-chan amqp.Delivery)) {
	return func (ch *amqp.Channel) (err error, msgs <-chan amqp.Delivery) {
		count := 0
		for{
			_, err = ch.QueueDeclare(args.Qname,args.Durable, args.AutoDelete, args.Exclusive, args.NoWait, args.Args)
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
			err = ch.QueueBind(args.Qname, args.RoutingKey, args.Exchange, args.NoWait, args.Args)
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
			err = ch.Qos(args.PrefetchCount, args.PrefetchSize, args.Global)
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
			msgs, err = ch.Consume(args.Qname, args.Consumer, args.AutoAck, args.Exclusive, args.NoLocal, args.NoWait, args.Args)
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
}