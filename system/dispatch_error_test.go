package system

import (
	"testing"
	"shangwoa.com/event"
)

func TestSystemDomain_OnError(t *testing.T) {
	name := "test"
	msg := "hello"
	GetSystemDomainInstance().EventDispatcher.AddEventListener(name, func(event *event.Event) {
		fmt.Println("TestSystemDomain_OnError", event)
	})
	GetSystemDomainInstance().OnError(errors.New(msg), name)
}
