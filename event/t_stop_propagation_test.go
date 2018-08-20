package event

import (
	"fmt"
	"testing"
)

func propagationHandler(event *Event) {
	fmt.Println("never print")
}
func printHandler(event *Event) {
	event.StopPropagation()
	fmt.Println("next stoped")
}

func TestEvent_StopPropagation(t *testing.T) {
	eventType := "eventType"
	dispatcher := NewEventDispatcher()

	//listener := NewListener(printHandler)
	dispatcher.AddEventListener(eventType, printHandler)

	//listener = NewListener(propagationHandler)
	dispatcher.AddEventListener(eventType,printHandler)

	event := NewEvent(eventType, "data", false)
	dispatcher.Dispatch(event)
}