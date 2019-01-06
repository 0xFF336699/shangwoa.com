package redis2

import (
	"fmt"
	"github.com/go-redis/redis"
)

var clients = map[int]*redis.Client{}
var Client *redis.Client
func CreateClient(addr, pw string, db int) (err error, client *redis.Client)  {
	if client, ok := clients[db]; ok{
		return nil, client
	}
	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw, // no password set
		DB:       db,  // use default DB
	})
	pong, err := client.Ping().Result()
	if err != nil{
		s := err.Error()
		panic(s)
	}
	fmt.Println("redis pong", pong)
	clients[db] = client
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