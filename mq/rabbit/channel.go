package rabbit

import (
	"github.com/streadway/amqp"
	"time"
)

type Channel struct{
	RetryMaxCount int
	Ch *amqp.Channel
	Connection *Connection
}

func (this *Channel) ReinitConnection()  {
	if this.Ch != nil{
		this.Ch.Close()
		this.Ch = nil
	}
}
func (this *Channel) Channel() (err error, ch *amqp.Channel) {
	err, conn := this.Connection.Connect()
	if err != nil{
		return
	}
	count := 0
	for{
		ch, err = conn.Channel()

		if err == nil{
			this.Ch = ch
			receiver := make(chan *amqp.Error)
			receiver = ch.NotifyClose(receiver)
			go func() {
				for {
					select{
					case <-receiver:
						this.Ch = nil
						return
					}
				}
			}()
			return
		}
		count ++
		if count > this.RetryMaxCount{
			return
		}
		time.Sleep(time.Second * 60 * time.Duration(count))

	}
	return
}