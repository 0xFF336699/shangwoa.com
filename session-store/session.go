package session_store

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/rbcervilla/redisstore"
	"github.com/gorilla/sessions"
)
var stores = map[string]*redisstore.RedisStore{}
func CreateStore(alias, keyPrefix string, client *redis.Client, opts *sessions.Options ) (err error) {
	store, err := redisstore.NewRedisStore(client)
	if err != nil {
		fmt.Println("failed to create redis store: ", err.Error())
		return
	}
	store.KeyPrefix(keyPrefix)
	store.Options(*opts)
	stores[alias] = store
	return
}

func MustGetStore(alias string) *redisstore.RedisStore  {
	if store, ok := stores[alias]; ok{
		return store
	}
	panic("no store name is:" + alias);
}