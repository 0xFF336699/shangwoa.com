package system

import (
	"shangwoa.com/event"
)
func init() {
	domain = &SystemDomain{
		EventDispatcher:event.NewEventDispatcher(),
	}
}
