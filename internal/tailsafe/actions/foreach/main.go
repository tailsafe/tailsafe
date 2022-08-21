package foreachAction

import (
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"reflect"
)

type Config struct {
	Use           any                   `json:"use"`
	ActionSetters tailsafe.ActionSetter `json:"action-setters"`
}

type ForeachAction struct {
	tailsafe.DataInterface
	tailsafe.StepInterface
	Config *Config

	global tailsafe.DataInterface
	result any
}

func (fa *ForeachAction) Configure() (err tailsafe.ErrActionInterface) {
	return
}

func (fa *ForeachAction) Execute() (err tailsafe.ErrActionInterface) {
	str, ok := fa.Config.Use.(string)
	if ok {
		fa.Config.Use = fa.Resolve(str, fa.GetAll())
	}

	rf := reflect.ValueOf(fa.Config.Use)
	if !rf.IsValid() {
		return
	}

	payload := tailsafe.NewPayload()

	// Inject current data into the payload
	for k, v := range fa.GetAll() {
		payload.Set(k, v)
	}

	if rf.Kind() != reflect.Slice {
		return
	}
	// iterate over the slice
	for i := 0; i < rf.Len(); i++ {

		if fa.Config.ActionSetters.Key != "" {
			payload.Set(fa.Config.ActionSetters.Key, i)
		}

		if fa.Config.ActionSetters.Value != "" {
			payload.Set(fa.Config.ActionSetters.Value, rf.Index(i).Interface())
		}

		err := fa.Next(payload)
		if err != nil {
			return err
		}
	}
	return
}

func (fa *ForeachAction) GetResult() interface{} {
	return fa.result
}

func (fa *ForeachAction) GetConfig() interface{} {
	return fa.Config
}

func (fa *ForeachAction) SetPayload(data tailsafe.DataInterface) {
	fa.DataInterface = data
}

func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(ForeachAction)
	p.StepInterface = step
	p.Config = new(Config)
	return p
}
