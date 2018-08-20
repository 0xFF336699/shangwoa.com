package event


var listenerIndex = 0
type Listener struct {
	handler func(event *Event)
	index int
}

func NewListener(handler func(event *Event)) (l *Listener) {
	l = &Listener{
		handler:handler,
		index:listenerIndex,
	}
	listenerIndex ++
	return
}
