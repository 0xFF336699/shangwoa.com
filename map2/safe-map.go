package map2

import "sync"

type IntIntMap struct {
	sync.RWMutex
	Map map[int]int
}


func NewIntIntMap() *IntIntMap {
	sm := new(IntIntMap)
	sm.Map = make(map[int]int)
	return sm
}

func (sm *IntIntMap)Get(key int)(int, bool) {
	sm.RLock()
	value, ok := sm.Map[key]
	sm.RUnlock()
	return value, ok
}

func (sm *IntIntMap)Set(key int, value int) {
	sm.Lock()
	sm.Map[key] = value
	sm.Unlock()
}

func (sm * IntIntMap)Delete(key int) {
	sm.Lock()
	delete(sm.Map, key)
	sm.Unlock()
}


type StringStringMap struct {
	sync.RWMutex
	Map map[string]string
}


func NewStringStringMap() *StringStringMap {
	sm := new(StringStringMap)
	sm.Map = make(map[string]string)
	return sm
}

func (sm *StringStringMap)Get(key string)(string, bool) {
	sm.RLock()
	value, ok := sm.Map[key]
	sm.RUnlock()
	return value, ok
}

func (sm *StringStringMap)Set(key string, value string) {
	sm.Lock()
	sm.Map[key] = value
	sm.Unlock()
}

func (sm * StringStringMap)Delete(key string) {
	sm.Lock()
	delete(sm.Map, key)
	sm.Unlock()
}


type StringIntMap struct {
	sync.RWMutex
	Map map[string]int
}


func NewStringIntMap() *StringIntMap {
	sm := new(StringIntMap)
	sm.Map = make(map[string]int)
	return sm
}

func (sm *StringIntMap)Get(key string)(int, bool) {
	sm.RLock()
	value, ok := sm.Map[key]
	sm.RUnlock()
	return value, ok
}

func (sm *StringIntMap)Set(key string, value int) {
	sm.Lock()
	sm.Map[key] = value
	sm.Unlock()
}

func (sm * StringIntMap)Delete(key string) {
	sm.Lock()
	delete(sm.Map, key)
	sm.Unlock()
}

func (sm *StringIntMap) Range() (keys []string, values []int) {

	return nil,nil
}
