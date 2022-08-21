package findAction

import (
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
)

type Config struct {
	Use           any                   `json:"use"`
	ActionGetters tailsafe.ActionGetter `json:"action-getters"`
}

type FindAction struct {
	tailsafe.DataInterface
	tailsafe.StepInterface
	Config *Config

	global tailsafe.DataInterface
	result any
}

func (fa *FindAction) Configure() (err tailsafe.ErrActionInterface) {
	return
}

func (fa *FindAction) Execute() (err tailsafe.ErrActionInterface) {
	return
}

func (fa *FindAction) GetResult() interface{} {
	return fa.result
}

func (fa *FindAction) GetConfig() interface{} {
	return fa.Config
}

func (fa *FindAction) SetPayload(data tailsafe.DataInterface) {
	fa.DataInterface = data
}

func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(FindAction)
	p.StepInterface = step
	p.Config = new(Config)
	return p
}
