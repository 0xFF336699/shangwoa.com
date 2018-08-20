package rabbitmq

import (
	"errors"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

var EmptyBodyError = errors.New("body is empty")
var RetryAlarmError = errors.New("retry alarm")

func waitingRetryListen(exec func(d *amqp.Delivery) error, onError func(qname string, err error, extra string) bool, url string, qd *QueueDeclare, qb *QueueBind, qos *Qos, consume *Consume, zombieTriggerTime time.Duration, retryCount int, retryAlarmCount int) {
	time.Sleep(5 * time.Second)
	retryListen(exec, onError, url, qd, qb, qos, consume, zombieTriggerTime, retryCount, retryAlarmCount)
}

func retryListen(exec func(d *amqp.Delivery) error, onError func(qname string, err error, extra string) bool, url string, qd *QueueDeclare, qb *QueueBind, qos *Qos, consume *Consume, zombieTriggerTime time.Duration, retryCount int, retryAlarmCount int) {
	retryCount++
	if retryCount == retryAlarmCount {
		stop := onError(qd.Queue, RetryAlarmError, "")
		if stop {
			return
		}
	}
	Listen(exec, onError, url, qd, qb, qos, consume, zombieTriggerTime, retryCount, retryAlarmCount)
}

func Listen(exec func(d *amqp.Delivery) error, onError func(qname string, err error, extra string) bool, url string, qd *QueueDeclare, qb *QueueBind, qos *Qos, consume *Consume, zombieTriggerTime time.Duration, retryCount int, retryAlarmCount int) {
	return 
	ch, err := GetChannelByDefaultRetries(url)
	if err != nil {
		go waitingRetryListen(exec, onError, url, qd, qb, qos, consume, zombieTriggerTime, retryCount, retryAlarmCount)
		return
	}
	defer ch.Close()
	_, err = ch.QueueDeclare(
		qd.Queue,      // name
		qd.Durable,    // durable
		qd.AutoDelete, // delete when unused
		qd.Exclusive,  // exclusive
		qd.NoWait,     // no-wait
		qd.Arguments,  // arguments
	)

	if err != nil {
		go waitingRetryListen(exec, onError, url, qd, qb, qos, consume, zombieTriggerTime, retryCount, retryAlarmCount)
		return
	}
	err = ch.QueueBind(
		qb.Queue,      // queue name
		qb.RoutingKey, // routing key
		qb.Exchange,   // exchange
		qb.NoWait,
		qb.Arguments)

	if err != nil {
		go waitingRetryListen(exec, onError, url, qd, qb, qos, consume, zombieTriggerTime, retryCount, retryAlarmCount)
		return
	}
	err = ch.Qos(
		qos.PrefetchCount, // prefetch count
		qos.PrefetchSize,  // prefetch size
		qos.Global,        // global
	)

	if err != nil {
		go waitingRetryListen(exec, onError, url, qd, qb, qos, consume, zombieTriggerTime, retryCount, retryAlarmCount)
		return
	}
	fmt.Println("listening is", qd.Queue)
	msgs, err := ch.Consume(
		consume.Queue,       // queue
		consume.ConsumerTag, // consumer
		consume.NoAck,       // auto-ack
		consume.Exclusive,   // exclusive
		consume.NoLocal,     // no-local
		consume.NoWait,      // no-wait
		consume.Arguments,   // args
	)

	if err != nil {
		go waitingRetryListen(exec, onError, url, qd, qb, qos, consume, zombieTriggerTime, retryCount, retryAlarmCount)
		return
	}
	//ticker := time.NewTicker(zombieTriggerTime * time.Second)
	//isGone := false
	//forever := make(chan bool)
	//breakMsgs := make(chan bool)
	//
	//go func() {
	//	fmt.Println("listen msgs")
	//	retryCount = 0
	//	for {
	//		select {
	//		case d := <-msgs:
	//			ticker.Stop()
	//			if len(d.Body) > 0 && isGone == false {
	//				err := exec(&d)
	//				if err != nil {
	//					stopedOnError = onError(qd.Queue, err, string(d.Body))
	//				}
	//			} else if len(d.Body) == 0 {
	//
	//			}
	//
	//			if stopedOnError {
	//				// 这样会导致后续的也不再执行，所以应该有启动提醒程序猿立即解决机制
	//				isGone = true
	//				forever <- false
	//				goto END
	//			}
	//			d.Ack(false)
	//			ticker = time.NewTicker(zombieTriggerTime * time.Second)
	//		case <-breakMsgs:
	//			isGone = true
	//			goto END
	//		}
	//	}
	//END:
	//}()
	//
	//go func() {
	//	for {
	//		select {
	//		case <-ticker.C:
	//			if !stopedOnError {
	//				go retryListen(exec, onError, url, qd, qb, qos, consume, zombieTriggerTime, retryCount, retryAlarmCount)
	//			}
	//			fmt.Println("time over")
	//			ticker.Stop()
	//			isGone = true
	//			breakMsgs <- true
	//			forever <- false
	//			goto END
	//		}
	//	}
	//END:
	//}()

	ticker := time.NewTicker(zombieTriggerTime * time.Second)
	isGone := false
	forever := make(chan bool)
	isActive := false
	activeCh := make(chan bool)
	breakMsgs := make(chan bool)
	var stopedOnError bool
	go func() {
		retryCount = 0
		for {
			select {
			case d := <-msgs:
				ticker.Stop()
				if len(d.Body) > 0 {
					err := exec(&d)
					if err != nil {
						stopedOnError = onError(qd.Queue, err, string(d.Body))
					}
				} else {
					stopedOnError = onError(qd.Queue, EmptyBodyError, "body is empty")

				}
				if stopedOnError {
					close(breakMsgs)
					// 这样会导致后续的也不再执行，所以应该有启动提醒程序猿立即解决机制
					goto END
				}
				ticker = time.NewTicker(zombieTriggerTime * time.Second)
				d.Ack(false)
			case <-breakMsgs:
				close(breakMsgs)
				ticker.Stop()
				goto END
			}
		}
	END:
	}()

	go func() {
		for {
			select {
			case bl := <-activeCh:
				isActive = bl
			case <-ticker.C:
				if !isActive {
					if !isGone {
						isGone = true
						breakMsgs <- true
						ch.Close()
						ticker.Stop()
						if !stopedOnError {
							go retryListen(exec, onError, url, qd, qb, qos, consume, zombieTriggerTime, retryCount, retryAlarmCount)
						}
						forever <- false
						goto END
					}
				} else {
					activeCh <- false
				}
			}
		}
	END:
	}()
	<-forever
	msgs = nil
	close(activeCh)
	close(forever)
}
