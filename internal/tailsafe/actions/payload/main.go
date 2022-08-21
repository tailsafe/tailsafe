package payloadAction

import (
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
)

type Config struct {
}

type Payload struct {
	tailsafe.StepInterface
	tailsafe.DataInterface

	Config map[string]any
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
	p.DataInterface = data
}
func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(Payload)
	p.StepInterface = step
	return p
}
