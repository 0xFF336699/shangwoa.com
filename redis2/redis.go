package redis2

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
)

var clients = map[string]*redis.Client{}
var Client *redis.Client
func CreateClient(addr, pw string, db int) (err error, client *redis.Client)  {
	name := addr + strconv.Itoa(db)
	if client, ok := clients[name]; ok{
		return nil, client
	}
	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw, // no password set
		DB:       db,  // use default DB
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil{
		s := err.Error()
		panic(s)
	}
	fmt.Println("redis pong", pong)
	clients[name] = client
	return err, client
}
func SetDefaultClient(addr, pw string, db int) (err error, client *redis.Client) {
	err, client = CreateClient(addr, pw, db)
	if err != nil{
		return
	}
	Client = client
	return
}

func Create(opt *redis.Options)(err error, client *redis.Client)  {
	client = redis.NewClient(opt)
	_, err = client.Ping(context.Background()).Result()
	if err != nil{
		return
	}
	return
}
func GetClient(opt *redis.Options) (err error, client *redis.Client) {
	name := opt.Addr + strconv.Itoa(opt.DB)
	client, ok := clients[name]
	if !ok{
		err, client = Create(opt)
		if err != nil{
			return
		}
		clients[name] = client
	}
	return
}