package tailsafe

type Logger interface {
	Log(LoggerPayload)
}
type LoggerPayload interface {
	GetNamespace() string
	GetLevel() int
	GetMessage() string
	GetArgs() []any

	SetNamespace(string) LoggerPayload
	SetLevel(int) LoggerPayload
	SetMessage(string) LoggerPayload
	SetArgs(...any) LoggerPayload
}
