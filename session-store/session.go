package session_store

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/sessions"
	"github.com/rbcervilla/redisstore"
)
var stores = map[string]*redisstore.RedisStore{}
func CreateStore(alias, keyPrefix string, client *redis.Client, opts *sessions.Options ) (err error) {
	store, err := redisstore.NewRedisStore(context.Background(), client)
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