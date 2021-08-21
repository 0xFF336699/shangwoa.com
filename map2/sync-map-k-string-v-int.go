package map2

import (

	"sync"
)
func SyncMapLoadKStringVInt(m sync.Map, key string)(v int, ok bool)  {
	i, ok := m.Load(key)
	if !ok{
		return
	}
	v = i.(int)
	return
}

func SyncMapLoadOrStoreKStringVInt(m sync.Map, key string, value int)(v int, ok bool)  {
	i, ok := m.LoadOrStore(key, value)
	if !ok{
		return
	}
	v = i.(int)
	return
}

func SyncMapRangeKStringVInt(m sync.Map, f func(key string, value int) bool )  {
	m.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(int)
		return f(k, v)
	})
	return
}

