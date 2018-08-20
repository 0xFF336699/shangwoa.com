package rabbitmq

import "github.com/streadway/amqp"

type MqConfVO struct {
	MqName	string	`json:"mq_name"`
	Queue QueueVO `json:"queue"`
	Exchange	ExchangeVO	`json:"exchange"`
	Publish	PublishVO	`json:"publish"`
	Publishing	*amqp.Publishing	`json:"publishing"`
}

type QueueVO struct{
	Name string `json:"name"`
	Durable bool 	`json:"durable"`
	AutoDelete bool `json:"auto_delete"`
	Exclusive bool `json:"exclusive"`
	NoWait bool `json:"no_wait"`

}

type ExchangeVO struct {
	NoWait	bool	`json:"no_wait"`
	Name	string	`json:"name"`
	Kind	string	`json:"kind"`
	Durable	bool	`json:"durable"`
	AutoDelete	bool	`json:"auto_delete"`
	Internal	bool	`json:"internal"`
}

type PublishVO struct {
	Immediate	bool	`json:"immediate"`
	Exchange	string	`json:"exchange"`
	RoutingKey	string	`json:"routing_key"`
	Mandatory	bool	`json:"mandatory"`
}

