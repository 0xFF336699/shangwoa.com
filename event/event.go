package event

type IEventHander interface {

}

type Event struct {
	eventType string
	Data interface{}
	cancelable bool
	target *EventDispatcher
	stoped bool
}

func NewEvent(eventType string, data interface{}, cancelable bool) *Event {
	return &Event{
		eventType:eventType,
		Data:data,
		cancelable:cancelable,
	}
}
func (this *Event) GetType() string {
	return this.eventType
}
func (this *Event) StopPropagation()  {
	if this.cancelable{
		this.stoped = true
	}
}

func (this *Event) GetTarget() *EventDispatcher  {
	return this.target
}
