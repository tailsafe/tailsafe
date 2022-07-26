package tailsafe

import "sync"

type payload struct {
	sync.Mutex
	data map[string]interface{}
}

func NewPayload() *payload {
	return &payload{
		data: make(map[string]interface{}),
	}
}

func (d *payload) Get(key string) interface{} {
	d.Lock()
	defer d.Unlock()

	return d.data[key]
}

func (d *payload) Set(key string, value interface{}) {
	d.Lock()
	defer d.Unlock()

	d.data[key] = value
}

func (d *payload) GetAll() map[string]interface{} {
	return d.data
}
