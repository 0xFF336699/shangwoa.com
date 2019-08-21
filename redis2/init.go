package redis2

import "github.com/go-redis/redis"

func init() {
	clients = map[string]*redis.Client{}
}
