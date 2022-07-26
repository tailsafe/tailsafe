package payloadAction

import (
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
)

type Config struct {
}

type Payload struct {
	tailsafe.StepInterface
	Config  map[string]any `json:"year"`
	Payload tailsafe.DataInterface
	_Data   map[string]interface{}
}

func (p *Payload) Configure() (err tailsafe.ErrActionInterface) {
	return
}

func (p *Payload) GetResult() any {
	return p.Config
}

func (p *Payload) Execute() (err tailsafe.ErrActionInterface) {
	return
}
func (p *Payload) GetConfig() any {
	if p.Config == nil {
		p.Config = make(map[string]any)
	}
	return &p.Config
}
func (p *Payload) SetPayload(data tailsafe.DataInterface) {
	p.Payload = data
}
func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(Payload)
	p.StepInterface = step
	return p
}
