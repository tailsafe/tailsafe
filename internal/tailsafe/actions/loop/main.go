package loop

import (
	"github.com/pkg/errors"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
)

type Config struct {
	Use string `json:"use"`
}
type Loop struct {
	tailsafe.DataInterface
	tailsafe.StepInterface
	Config *Config
}

func (e *Loop) Configure() (err tailsafe.ErrActionInterface) {
	if e.Config == nil {
		return tailsafe.CatchStackTrace(e.GetContext(), errors.New("Loop: Config is nil"))
	}
	return
}

func (e *Loop) GetResult() any {
	return nil
}

func (e *Loop) Execute() (err tailsafe.ErrActionInterface) {
	v := e.Get(e.Config.Use)
	for _, _ = range v.([]interface{}) {
		// set current data
		// e.SetCurrent(e.GetKey(), value)

		// call next action
		/*err = e.Next(nil)
		if err != nil {
			var Err *tailsafe.ErrContinue
			if !errors.As(err.GetOriginal(), &Err) {
				return tailsafe.CatchStackTrace(e.GetContext(), err)
			}
			err = nil
			continue
		}*/
		/*if e.Plugin() == nil {
			continue
		}
		if e.Plugin().Data() == nil {
			continue
		}
		data := e.Plugin().Data().(map[string]interface{})
		v, ok := data["current"]
		if !ok {
			continue
		}
		v.([]map[string]interface{})[k] = value*/
	}
	return
}
func (e *Loop) Data() interface{} {
	return nil
}
func (e *Loop) GetConfig() interface{} {
	return e.Config
}
func (e *Loop) SetPayload(data tailsafe.DataInterface) {
	e.DataInterface = data
}
func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(Loop)
	p.StepInterface = step

	return p
}
