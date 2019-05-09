package rabbit

import (
	"fmt"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

type connMap struct {
	sync.Mutex
	m sync.Map
}
var m *connMap

type Connection struct{
	Url string
	RetryMaxCount int
	Conn *amqp.Connection
}

func (this *Connection)Restart() (err error) {
	fmt.Println("connection restart")
	err = this.Conn.Close()
	if err != nil{
		this.Connect()
	}
	fmt.Println("connection restarted")
	return
}
func (this *Connection) Connect()(err error, conn *amqp.Connection)  {
	count := 0
	for{
		err, conn = GetConnection(this.Url)
		if err == nil{
			this.Conn = conn
			receiver := make(chan *amqp.Error)
			receiver = conn.NotifyClose(receiver)
			go func() {
				for{
					select {
					case _, ok:=<-receiver:
						_ = ok //回头得想想ok=false应该怎么处理
						//this.Conn = nil
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
		time.Sleep(time.Second * 10 * time.Duration(count))
	}
	return
}

func (this *Connection) Destroy() {

}
func GetConnection(url string) (err error, conn *amqp.Connection) {
	key := getMapStoreKeyByUrl(url)
	c, ok := m.m.Load(key)
	if ok{
		conn, ok = c.(*amqp.Connection)
		//conn.Close()
		if !ok{
			// 可能需要注销这个connection，现在没时间，以后再做处理
		}
		return
	}
	conn, err = amqp.Dial(url)
	if err != nil{
		return
	}
	m.m.Store(key, conn)
	go func() {
		receiver := make(chan *amqp.Error)
		conn.NotifyClose(receiver)
		for{
			select {
			case <-receiver:
				onConnectionClosed(key, url, conn)
				return
			}
		}
	}()
	return
}

func onConnectionClosed(key, url string, conn *amqp.Connection)  {
	m.m.Delete(key)
}

// 添加这个方法的目的是为了防止和以前的方法碰撞
func getMapStoreKeyByUrl(url string)(string)  {
	return "key:" + url
}
