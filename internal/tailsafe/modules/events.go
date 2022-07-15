package modules

import "github.com/tailsafe/tailsafe/pkg/tailsafe"

type Events struct {
	callback func(tailsafe.EventInterface)
}

type Event struct {
	Key  string
	Body []any
}

func (e *Event) SetKey(key string) {
	e.Key = key
}

func (e *Event) SetBody(body ...any) {
	e.Body = body
}

func (e *Event) GetKey() string {
	return e.Key
}

func (e *Event) GetBody() []any {
	return e.Body
}

func NewEvent(key string, body ...any) *Event {
	return &Event{
		Key:  key,
		Body: body,
	}
}

var EventsInstance *Events

func init() {
	EventsInstance = &Events{}
}

func GetEventsModule() *Events {
	return EventsInstance
}

func (e *Events) Subscribe(callback func(eventInterface tailsafe.EventInterface)) *Events {
	e.callback = callback
	return e
}

func (e *Events) Trigger(event tailsafe.EventInterface) {
	e.callback(event)
}
