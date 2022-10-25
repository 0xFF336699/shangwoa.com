package map2

import (
	"sync"
)
func KStringVInt64SyncMapLoad(m sync.Map, key string)(v int64, ok bool)  {
	i, ok := m.Load(key)
	if !ok{
		return
	}
	v = i.(int64)
	return
}

func KStringVInt64SyncMapLoadOrStore(m sync.Map, key string, value int64)(v int64, ok bool)  {
	i, ok := m.LoadOrStore(key, value)
	if !ok{
		return
	}
	v = i.(int64)
	return
}

func KStringVInt64SyncMapRange(m sync.Map, f func(key string, value int64) bool )  {
	m.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(int64)
		return f(k, v)
	})
	return
}

