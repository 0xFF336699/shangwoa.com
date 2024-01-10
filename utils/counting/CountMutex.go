package counting

import (
	"shangwoa.com/broadcast"
	"sync"
)

type CountMutex struct{
	M sync.Mutex
	Count int
	Down bool
	BroadCast *broadcast.Broadcast
}

func NewCountMutex(c int) *CountMutex  {
	return &CountMutex{
		M:         sync.Mutex{},
		Count:     c,
		Down:      false,
		BroadCast: broadcast.NewBroadcaster(),
	}
}
// 在某些时候，例如第一次启动失败后重新尝试启动，这时会需要对计数器进行重置
func (this *CountMutex)Reset(c int) {
	this.M.Lock()
	this.Count = c
	this.Down = false
	this.M.Unlock()
}
func (this *CountMutex)Add(c int)  {
	this.M.Lock()
	this.Count += c
	this.M.Unlock()
}
func (this *CountMutex)Minus(c int)  {
	this.M.Lock()
	defer this.M.Unlock()
	if this.Down {
		return
	}
	this.Count -= c
	if this.Count <= 0{
		if this.Down == false{
			this.Down = true
			this.BroadCast.Broadcast(true)
		}
	}
}