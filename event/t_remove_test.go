package event

import (
	"testing"
	"fmt"
)

func removeHandler(event *Event) {
	fmt.Println("never print")
}
func TestEventDispatcher_RemoveListener(t *testing.T) {
	eventType := "eventType"
	//listener := NewListener(removeHandler)
	dispatcher := NewEventDispatcher()
	dispatcher.AddEventListener(eventType, removeHandler)

	removed1:=dispatcher.RemoveListener(eventType, removeHandler)

	event := NewEvent(eventType, "data", false)
	dispatcher.Dispatch(event)

	removed2 := dispatcher.RemoveListener(eventType, removeHandler)
	fmt.Println("removed Eventlistener", removed1, removed2)
}