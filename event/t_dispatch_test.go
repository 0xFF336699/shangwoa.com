package event

import (
	"testing"
	"fmt"
)

func dispatchHandler(event *Event) {
	fmt.Println("i'm listening",event)
}
func TestEventDispatcher_Dispatch(t *testing.T) {

	eventType := "eventType"
	//listener := NewListener(dispatchHandler)
	dispatcher := NewEventDispatcher()
	dispatcher.AddEventListener(eventType, dispatchHandler)
	event := NewEvent(eventType, "data", false)
	dispatcher.Dispatch(event)
}


func eventHandler(event *Event) {
	fmt.Println("event is %+v", event)
}