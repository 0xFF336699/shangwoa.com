package map2

import (
	"sync"
	"shangwoa.com/time2"
)
func KIntVJsonTimeSyncMapLoad(m sync.Map, key int)(v time2.JsonTime, ok bool)  {
	i, ok := m.Load(key)
	if !ok{
		return
	}
	v = i.(time2.JsonTime)
	return
}

func KIntVJsonTimeSyncMapLoadOrStore(m sync.Map, key int, value time2.JsonTime)(v time2.JsonTime, ok bool)  {
	i, ok := m.LoadOrStore(key, value)
	if !ok{
		return
	}
	v = i.(time2.JsonTime)
	return
}

func KIntVJsonTimeSyncMapRange(m sync.Map, f func(key int, value time2.JsonTime) bool )  {
	m.Range(func(key, value interface{}) bool {
		k := key.(int)
		v := value.(time2.JsonTime)
		return f(k, v)
	})
	return
}

