package rabbitmq2

import (
	"github.com/streadway/amqp"
	"sync"
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