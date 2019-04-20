package rabbitmq

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/streadway/amqp"
	"shangwoa.com/system"
	"shangwoa.com/utils/retry"
)

type connMap struct {
	sync.Mutex
	m sync.Map
}

const (
	RabbitmqConnectionCreateError = "RabbitmqConnectionCreateError"
)

type connection struct {
	sync.Mutex
	m sync.Map
	//m map[string] *channel
	conn *amqp.Connection
}

type channel struct {
	sync.Mutex
	m sync.Map
	//m map[string] *QueueDeclare
	ch *amqp.Channel
}

type QueueDeclare struct {
	reserved1  uint16
	Queue      string
	Passive    bool
	Durable    bool
	Exclusive  bool
	AutoDelete bool
	NoWait     bool
	Arguments  amqp.Table
}
type QueueBind struct {
	reserved1  uint16
	Queue      string
	Exchange   string
	RoutingKey string
	NoWait     bool
	Arguments  amqp.Table
}
type Qos struct {
	PrefetchSize  int
	PrefetchCount int
	Global        bool
}

type Consume struct {
	reserved1   uint16
	Queue       string
	ConsumerTag string
	NoLocal     bool
	NoAck       bool
	Exclusive   bool
	NoWait      bool
	Arguments   amqp.Table
}
type ConnectionConf struct {
	FirstLevelMaxRetryCount             int
	FirstLevelMaxRetryCountWaitingTime  int
	SecondLevelMaxRetryCount            int
	SecondLevelMaxRetryCountWaitingTime int
}
type ChannelConf struct {
	FirstLevelMaxRetryCount             int
	FirstLevelMaxRetryCountWaitingTime  int
	SecondLevelMaxRetryCount            int
	SecondLevelMaxRetryCountWaitingTime int
}
type connectionLoadInfo struct {
	firstLevelRetryCount  int
	errors                []*error
	secondLevelRetryCount int
}

var m *connMap

var DefaultConnConf *ConnectionConf

//func createConnection(name string, connConf *ConnectionConf, loadInfo *connectionLoadInfo) (*connection, error) {
//	m.Lock()
//	defer m.Unlock()
//	return createConn(name, connConf, loadInfo)
//}
//
//func createConn(name string, connConf *ConnectionConf, loadInfo *connectionLoadInfo) (*connection, error) {
//	conn, err := amqp.Dial(name)
//	if err != nil{
//		if loadInfo.errors == nil{
//			loadInfo.errors = []*error{}
//		}
//		loadInfo.errors = append(loadInfo.errors, &err)
//		if loadInfo.firstLevelRetryCount < connConf.FirstLevelMaxRetryCount{
//			loadInfo.firstLevelRetryCount ++
//			time.Sleep(time.Duration(connConf.SecondLevelMaxRetryCountWaitingTime) * time.Second)
//			return createConn(name, connConf, loadInfo)
//		}else if loadInfo.secondLevelRetryCount < connConf.SecondLevelMaxRetryCount{
//			loadInfo.secondLevelRetryCount ++
//			time.Sleep(time.Duration(connConf.SecondLevelMaxRetryCountWaitingTime) * time.Second)
//			return createConn(name, connConf, loadInfo)
//		}else{
//			e:=&system.Error{
//				Type:RabbitmqConnectionCreateError,
//				Err:err,
//				SubData:loadInfo,
//			}
//			system.OnError(e)
//			return nil, e
//		}
//	}
//	return &connection{conn:conn}, err
//}
//func getConnection(name string, connConf *ConnectionConf, loadInfo *connectionLoadInfo) (*connection, error) {
//	conn, ok := m.m.Load(name)
//	if !ok{
//		conn, err:= createConnection(name, connConf, loadInfo)
//		if err != nil{
//			return nil, err
//		}
//		m.m.Store(name, conn)
//		return conn, err
//	}
//	return conn.(*connection), nil
//}

//func GetConnByDefaultConf(connName string) (*amqp.Connection, error) {
//	conn, err := getConnection(connName, DefaultConnConf, &connectionLoadInfo{})
//	if err != nil{
//		return nil,err
//	}
//	return conn.conn, err
//}

func GetConn(url string, retries *retry.Retries) (*amqp.Connection, error) {
	fmt.Println("Get conn url is", url, retries.Retry.Count)
	conn, err := getConn(url, retries)

	if err != nil {
		return nil, err
	}
	return conn.conn, err
}

func getConn(url string, retries *retry.Retries) (*connection, error) {
	conn, ok := m.m.Load(url)
	if !ok {
		conn, err := createconn(url, retries)
		if err != nil {
			return nil, err
		}
		m.m.Store(url, conn)
		return conn, err
	}
	fmt.Println("use old connection url is", url)
	return conn.(*connection), nil
}

func killConn(url string) (err error) {
	conn, ok := m.m.Load(url)
	if ok{
		c := conn.(*connection)
		if c != nil{
			if c.conn != nil{
				c.conn.Close()
			}
		}
	}
	return
}
func createconn(url string, retries *retry.Retries) (*connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		for retries.Retry != nil && retries.Retry.Count < retries.Retry.Max {
			retries.Retry.Errors = append(retries.Retry.Errors, &err)
			retries.Retry.Count++
			time.Sleep(retries.Retry.WaitingTime)
			conn, err = amqp.Dial(url)
			if err == nil {
				return &connection{conn: conn}, err
			}
			if retries.Retry.Count >= retries.Retry.Max {
				bl := retries.Next()
				if !bl {
					e := &system.Error{Type: "mq", Err: err, SubType: "retries.Next", SubData: retries, CodeLevel: 1}
					return nil, e
				}
			}
		}
	}
	return &connection{conn: conn}, err
}

func GetChannel(url string, connRetries *retry.Retries, channelRetries *retry.Retries) (*amqp.Channel, error) {
	conn, err := getConn(url, connRetries)
	if err != nil {
		return nil, err
	}
	ch, err := conn.conn.Channel()
	if err != nil {
		fmt.Println("GetChannel conn.con.Channel error", err.Error())
		conn, err = reconnectConn(url, connRetries)
		if err != nil {
			return nil, err
		}
		ch, err = conn.conn.Channel()
	}
	return ch, err
}

func GetChannelByDefaultRetries(url string) (*amqp.Channel, error) {
	conn := getConnRetries()
	ch := getChannelRetries()
	return GetChannel(url, &conn, &ch)
}

var defConnRetries *retry.Retries
var defChannelRetries *retry.Retries

func getConnRetries() retry.Retries {
	return *defConnRetries
}

func getChannelRetries() retry.Retries {
	return *defChannelRetries
}
func reconnectConn(url string, retries *retry.Retries) (*connection, error) {
	conn, err := createconn(url, retries)
	if err != nil {
		return nil, err
	}
	m.m.Store(url, conn)
	return conn, err
}

func init() {
	m = &connMap{}
	DefaultConnConf = &ConnectionConf{
		FirstLevelMaxRetryCount:             1,
		FirstLevelMaxRetryCountWaitingTime:  1,
		SecondLevelMaxRetryCount:            1,
		SecondLevelMaxRetryCountWaitingTime: 1,
	}

	defConnRetries = NewRetries(3, 1*time.Second, 2, 2*time.Second)
	defChannelRetries = NewRetries(3, 1*time.Second, 2, 2*time.Second)
}

func NewRetries(args ...interface{}) *retry.Retries {
	l := len(args)
	retries := &retry.Retries{
		List: list.New(),
	}
	for i := 0; i < l; i += 2 {
		d := args[i+1]
		r := &retry.Retry{
			Max:         args[i].(int),
			WaitingTime: d.(time.Duration),
		}
		retries.List.PushBack(r)
	}
	retries.Reset()
	return retries
}
