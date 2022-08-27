package data

import (
	"encoding/json"
	"github.com/tailsafe/tailsafe/internal/tailsafe/data/setter"
	"sync"
)

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

func (d *Payload) Set(key string, value any, override bool) {
	d.mx.Lock()
	defer d.mx.Unlock()

	v, err := json.Marshal(value)
	if err != nil {
		return
	}

	// force untyped data, yes why not 😱
	var slice any
	err = json.Unmarshal(v, &slice)
	if err != nil {
		return
	}

	err = setter.
		Set(d.data).
		SetOverride(override).
		Apply(key, slice)

	return
}

func (d *Payload) GetAll() map[string]any {
	return d.data
}
