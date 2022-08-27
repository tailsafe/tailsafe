package tailsafe

type DataInterface interface {
	Set(key string, value any, override bool)
	Get(key string) any

	GetAll() map[string]any
}
