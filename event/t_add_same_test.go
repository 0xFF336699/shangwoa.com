package event

import (
	"testing"
	"fmt"
)

type Func func(event *Event)


func addSameEventTypeHandler(event *Event) {
	fmt.Println("just one time")
}
func TestEventDispatcher_AddEventListener(t *testing.T) {
	eventType := "eventType"
	//listener := NewListener(addSameEventTypeHandler)
	dispatcher := NewEventDispatcher()
	dispatcher.AddEventListener(eventType, addSameEventTypeHandler)
	dispatcher.AddEventListener(eventType,addSameEventTypeHandler)

	event := NewEvent(eventType, "data", false)
	dispatcher.Dispatch(event)
}