package jsonEncodeAction

import (
	"encoding/json"
	"errors"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
)

type Config struct {
	Value any `yaml:"value"`
}

type JsonEncodeAction struct {
	tailsafe.DataInterface
	tailsafe.StepInterface
	Config *Config

	global tailsafe.DataInterface
	result string
}

func (pa *JsonEncodeAction) Configure() (err tailsafe.ErrActionInterface) {
	if pa.Config.Value == nil {
		return tailsafe.CatchStackTrace(pa.GetContext(), errors.New("JsonEncodeAction: Value cannot be empty"))
	}
	return
}

func (pa *JsonEncodeAction) Execute() tailsafe.ErrActionInterface {
	str, ok := pa.Config.Value.(string)
	if ok {
		pa.Config.Value = pa.Resolve(str, pa.GetAll())
	}
	res, err := json.Marshal(pa.Config.Value)
	if err != nil {
		return tailsafe.CatchStackTrace(pa.GetContext(), err)
	}

	pa.result = string(res)
	return nil
}
func (pa *JsonEncodeAction) GetResult() interface{} {
	return pa.result
}
func (pa *JsonEncodeAction) GetConfig() interface{} {
	return pa.Config
}
func (pa *JsonEncodeAction) SetPayload(data tailsafe.DataInterface) {
	pa.DataInterface = data
}
func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(JsonEncodeAction)
	p.StepInterface = step
	p.Config = new(Config)
	return p
}
