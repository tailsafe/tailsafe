package tailsafe

type DataInterface interface {
	Set(key string, value any)
	Get(key string) any
	GetAll() map[string]any
}
