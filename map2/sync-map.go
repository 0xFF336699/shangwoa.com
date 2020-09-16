package map2

import (
	"os"
	"sync"
)

type StringIFileSyncMap struct {
	Map sync.Map
}

func (this *StringIFileSyncMap) Load(key string)(v *os.File, ok bool)  {
	i, ok := this.Map.Load(key)
	if !ok{
		return
	}
	v = i.(*os.File)
	return
}

func (this *StringIFileSyncMap) LoadOrStore(key string, value *os.File)(v *os.File, ok bool)  {
	i, ok := this.Map.LoadOrStore(key, value)
	if !ok{
		return
	}
	v = i.(*os.File)
	return
}

func (this *StringIFileSyncMap) Range(f func(key string, value *os.File) bool )  {
	this.Map.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(*os.File)
		return f(k, v)
	})
	return
}
