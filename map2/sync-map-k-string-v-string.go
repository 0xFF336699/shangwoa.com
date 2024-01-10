package map2

import (
	"sync"
)
//如果要外部可访问，就把name大写，生成代码后，再把map的变量名重构成小写，这样就可以方法是外部可访问，变量不可访问了

func KStringVStringSyncMapLoad(m sync.Map, key string)(v string, ok bool)  {
	i, ok := m.Load(key)
	if !ok{
		return
	}
	v = i.(string)
	return
}

func KStringVStringSyncMapLoadOrStore(m sync.Map, key string, value string)(v string, ok bool)  {
	i, ok := m.LoadOrStore(key, value)
	if !ok{
		return
	}
	v = i.(string)
	return
}

func KStringVStringSyncMapRange(m sync.Map, f func(key string, value string) bool )  {
	m.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		return f(k, v)
	})
	return
}

