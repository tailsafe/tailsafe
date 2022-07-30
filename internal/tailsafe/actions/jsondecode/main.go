package jsonDecodeAction

import (
	"encoding/json"
	"errors"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"strings"
)

type Config struct {
	Value string `yaml:"value"`
}

type JsonDecodeAction struct {
	tailsafe.DataInterface
	tailsafe.StepInterface
	Config *Config

	global tailsafe.DataInterface
	result any
}

func (pa *JsonDecodeAction) Configure() (err tailsafe.ErrActionInterface) {
	if strings.TrimSpace(pa.Config.Value) == "" {
		return tailsafe.CatchStackTrace(pa.GetContext(), errors.New("JsonDecodeAction: Value cannot be empty"))
	}
	return
}

func (pa *JsonDecodeAction) Execute() tailsafe.ErrActionInterface {
	value, ok := pa.Resolve(pa.Config.Value, pa.GetAll()).(string)
	if !ok {
		return tailsafe.CatchStackTrace(pa.GetContext(), errors.New("JsonDecodeAction: Value need to be a string"))
	}

	err := json.Unmarshal([]byte(value), &pa.result)
	if err != nil {
		return tailsafe.CatchStackTrace(pa.GetContext(), err)
	}
	return nil
}
func (pa *JsonDecodeAction) GetResult() any {
	return pa.result
}
func (pa *JsonDecodeAction) GetConfig() any {
	return pa.Config
}
func (pa *JsonDecodeAction) SetPayload(data tailsafe.DataInterface) {
	pa.DataInterface = data
}
func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(JsonDecodeAction)
	p.StepInterface = step
	p.Config = new(Config)
	return p
}
