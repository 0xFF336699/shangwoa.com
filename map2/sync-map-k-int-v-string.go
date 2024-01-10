package map2

import (

	"sync"
)
var KIntVStringSyncMap sync.Map
func KIntVStringSyncMapLoad(m sync.Map, key int)(v string, ok bool)  {
	i, ok := m.Load(key)
	if !ok{
		return
	}
	v = i.(string)
	return
}

func KIntVStringSyncMapLoadOrStore(m sync.Map, key int, value string)(v string, ok bool)  {
	i, ok := m.LoadOrStore(key, value)
	if !ok{
		return
	}
	v = i.(string)
	return
}

func KIntVStringSyncMapRange(m sync.Map, f func(key int, value string) bool )  {
	m.Range(func(key, value interface{}) bool {
		k := key.(int)
		v := value.(string)
		return f(k, v)
	})
	return
}
