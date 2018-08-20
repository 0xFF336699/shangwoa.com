package system

import (
	"shangwoa.com/event"
)
type SystemDomain struct{
	EventDispatcher *event.EventDispatcher
}

func (this *SystemDomain)OnError(err error, name string)  {
	this.DispatchError(name, err)
}
func (this *SystemDomain) DispatchError(eventType string, err error) {
	evt:= event.NewEvent(eventType, err, false)
	this.EventDispatcher.Dispatch(evt)
}
var domain *SystemDomain
func GetSystemDomainInstance() *SystemDomain {
	return domain
}
