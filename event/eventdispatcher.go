package event

import (
	"container/list"
	"reflect"
)

type Handler func(event *Event)

type IEventDispatcher interface {
	AddEventListener(eventType string, handler Handler)
	RemoveListener(eventType string, handler Handler) bool
	Dispatch(event *Event) bool
	HasEventListener(eventType string) bool
} 


type EventDispatcher struct{
	types map[string]*list.List
}

func (this *EventDispatcher) AddEventListener(eventType string, handler Handler) {
	if _, ok := this.types[eventType]; ok == false{
		this.types[eventType] = list.New()
	}
	l := this.types[eventType]
	for node := l.Front(); node != nil; node = node.Next(){
		if reflect.ValueOf(node.Value) == reflect.ValueOf(handler){
			l.Remove(node)
		}
		//if f, ok := node.Value.(Handler); ok {
		//	if f == listener{
		//		l.Remove(node)
		//	}
		//}
	}
	l.PushBack(handler)
}

func (this *EventDispatcher) RemoveListener(eventType string, handler Handler) bool {
	if _, ok := this.types[eventType]; ok == false{
		return false
	}
	l := this.types[eventType]
	for node := l.Front(); node != nil; node = node.Next(){
		if reflect.ValueOf(node.Value) == reflect.ValueOf(handler){
			l.Remove(node)
			return true
		}
		//if f, ok := node.Value.(*Handler); ok {
		//	if f == listener{
		//		l.Remove(node)
		//		return true
		//	}
		//}
	}
	return false
}

func (this *EventDispatcher) Dispatch(event *Event) bool  {
	event.target = this
	eventType := event.eventType
	if _, ok := this.types[eventType]; ok == false{
		return false
	}
	l := this.types[eventType]
	for node := l.Front(); node != nil; node = node.Next(){
		if f, ok := node.Value.(Handler); ok {
			f(event)
			//f.handler(event)
			if event.stoped{
				return false
			}
		}
	}
	return true
}

func (this *EventDispatcher) HasEventListener(eventType string) bool {

	if _, ok := this.types[eventType]; ok == false{
		return false
	}
	l := this.types[eventType]
	return l.Len() > 0
	//for node := l.Front(); node != nil; node = node.Next(){
	//	if _, ok := node.Value.(*Handler); ok {
	//		return true
	//	}
	//}
	//return false
}

func NewEventDispatcher() ( *EventDispatcher) {
	dispatcher :=&EventDispatcher{
	}
	dispatcher.types = make(map[string]*list.List)
	return dispatcher
}