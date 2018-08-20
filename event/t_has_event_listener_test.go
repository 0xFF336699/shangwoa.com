package event

import (
	"testing"
	"fmt"
)

func TestEventDispatcher_HasEventListener(t *testing.T) {
	eventType := "eventType"
	dispatcher := NewEventDispatcher()
	hasListener := dispatcher.HasEventListener(eventType)
	fmt.Println("hasListener ", hasListener)
}
