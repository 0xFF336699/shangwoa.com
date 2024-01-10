package map2

import (
	"sync"
	"shangwoa.com/time2"
)
func KStringVJsonTimeSyncMapLoad(m sync.Map, key string)(v time2.JsonTime, ok bool)  {
	i, ok := m.Load(key)
	if !ok{
		return
	}
	v = i.(time2.JsonTime)
	return
}

func KStringVJsonTimeSyncMapLoadOrStore(m sync.Map, key string, value time2.JsonTime)(v time2.JsonTime, ok bool)  {
	i, ok := m.LoadOrStore(key, value)
	if !ok{
		return
	}
	v = i.(time2.JsonTime)
	return
}

func KStringVJsonTimeSyncMapRange(m sync.Map, f func(key string, value time2.JsonTime) bool )  {
	m.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(time2.JsonTime)
		return f(k, v)
	})
	return
}

