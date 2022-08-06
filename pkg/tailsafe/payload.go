package tailsafe

import "sync"

const THIS = "this"
const RETURN = "RETURN"

type payload struct {
	sync.Mutex
	data map[string]any
}

func NewPayload() *payload {
	return &payload{
		data: make(map[string]any),
	}
}

func (d *payload) Get(key string) any {
	d.Lock()
	defer d.Unlock()

	return d.data[key]
}

func (d *payload) Set(key string, value any) {
	d.Lock()
	defer d.Unlock()

	d.data[key] = value
}

func (d *payload) GetAll() map[string]any {
	return d.data
}
