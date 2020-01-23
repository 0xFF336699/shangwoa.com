package redisdb

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"math"
	"shangwoa.com/consul"
	"shangwoa.com/redis2"
	"time"
)
func InitRedis(kvs ...*consul.KVRedis) {
	for i := 0; i < len(kvs); i++{
		getRedisByKVRedis(kvs[i])
	}
}
type ClientWithConf struct{
	kvRedis *consul.KVRedis
	client *redis.Client
}
var clients = make(map[RedisName]*ClientWithConf)
var RedisConfNotFound = errors.New("redis config not found")
type RedisName string
func GetRedisByName(name RedisName) (err error, client *redis.Client)  {
	cc, ok := clients[name]
	if !ok{
		return RedisConfNotFound, nil
	}
	if cc.client == nil{
		err, client = getRedisByKVRedis(cc.kvRedis)
	}else{
		client = cc.client
	}
	return
}
func getRedisByKVRedis(kvRedis *consul.KVRedis) (err error, client *redis.Client) {
	if cc, ok := clients[RedisName(kvRedis.Name)]; ok{
		return nil, cc.client
	}

	opt := &redis.Options{
		Addr:     kvRedis.Addr,
		Password: kvRedis.PW, // no password set
		DB:       kvRedis.DB, // use default DB
	}
	err, client = createRedisClient(opt)
	if client != nil{
		clients[RedisName(kvRedis.Name)] = &ClientWithConf{
			kvRedis: kvRedis,
			client:  client,

		}
	}else{
		// panic
	}
	return
}
func createRedisClient(opt *redis.Options) (err error, client *redis.Client) {
	count := 0
	for {
		var err error
		err, client = redis2.GetClient(opt)
		if client != nil {
			_, err = client.Ping().Result()
			if err == nil{
				fmt.Println(fmt.Printf("redis connected %+v", opt))
				return nil, client
			}
		}
		count++
		fmt.Println("get redis err", err.Error())
		time.Sleep(time.Second * time.Duration(math.Min(20, float64(count))))
	}
	return
}