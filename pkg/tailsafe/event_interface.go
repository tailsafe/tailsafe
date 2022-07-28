package tailsafe

type EventsInterface interface {
	Trigger(event EventInterface)
}

type Event interface {
	SetBody(...any)
	SetKey(string)

	GetBody() []any
	GetKey() string
}
