package tailsafe

import "sync"

const THIS = "this"
const RETURN = "RETURN"

type Payload struct {
	mx   sync.Mutex
	data map[string]any
}

func NewPayload() *Payload {
	return &Payload{
		data: make(map[string]any),
	}
}

func (d *Payload) Get(key string) any {
	d.mx.Lock()
	defer d.mx.Unlock()

	return d.data[key]
}

func (d *Payload) Set(key string, value any) {
	d.mx.Lock()
	defer d.mx.Unlock()

	d.data[key] = value
}

func (d *Payload) GetAll() map[string]any {
	return d.data
}
