package storage

type IStorage interface {
	Set(table *interface{}, value interface{})
	Get(table *interface{})(value interface{}, err error)
	HSet(table *interface{}, value interface{})
	HGet(table *interface{}) (value interface{}, err error)
	HMSet(table *interface{}, ...interface{})
	HMGet(table *interface{}, ...interface{})(m map[interface{}]interface{}, err error)
}