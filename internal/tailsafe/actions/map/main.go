package mapAction

import (
	"errors"
	"fmt"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"strings"
)

type Config struct {
	Use   any    `json:"use"`
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type MapAction struct {
	tailsafe.DataInterface
	tailsafe.StepInterface
	Config *Config

	global tailsafe.DataInterface
	result map[string]any
}

func (ma *MapAction) Configure() (err tailsafe.ErrActionInterface) {
	if strings.TrimSpace(ma.Config.Key) == "" {
		return tailsafe.CatchStackTrace(ma.GetContext(), errors.New("MapAction: Key cannot be empty"))
	}
	str, ok := ma.Config.Use.(string)
	if !ok {
		return
	}

	v := ma.Resolve(str, ma.GetAll())
	mp, ok := v.(map[string]any)
	if !ok {
		return
	}

	ma.result = mp
	return
}

func (ma *MapAction) Execute() (err tailsafe.ErrActionInterface) {
	var value any
	str, ok := ma.Config.Value.(string)
	if ok {
		value = ma.Resolve(str, ma.GetAll())
	}

	ma.result[fmt.Sprintf("%v", ma.Resolve(ma.Config.Key, ma.GetAll()))] = value
	return
}

func (ma *MapAction) GetResult() interface{} {
	return ma.result
}

func (ma *MapAction) GetConfig() interface{} {
	return ma.Config
}

func (ma *MapAction) SetPayload(data tailsafe.DataInterface) {
	ma.DataInterface = data
}

func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(MapAction)
	p.StepInterface = step
	p.Config = new(Config)
	p.result = make(map[string]any)
	return p
}
