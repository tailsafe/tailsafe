package sortAction

import (
	"errors"
	"github.com/tailsafe/tailsafe/internal/tailsafe/data"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"golang.org/x/exp/slices"
)

type Config struct {
	ActionSetter tailsafe.ActionSetter `json:"action-setter"`
	ActionGetter tailsafe.ActionGetter `json:"action-getter"`

	Result struct {
		ActionGetter tailsafe.ActionGetter `json:"action-getter"`
	} `json:"result"`
}
type Sort struct {
	tailsafe.StepInterface
	tailsafe.DataInterface
	Config  *Config
	Payload any
}

func (s *Sort) Configure() (err tailsafe.ErrActionInterface) {
	if s.Config == nil {
		return tailsafe.CatchStackTrace(s.GetContext(), errors.New("sort: Config is nil"))
	}
	return
}

func (s *Sort) GetResult() any {
	return s.Payload
}

func (s *Sort) Execute() (err tailsafe.ErrActionInterface) {

	s.Payload = s.Resolve(s.Config.ActionGetter.Value, s.GetAll())

	slice, ok := s.Payload.([]interface{})
	if !ok {
		return tailsafe.CatchStackTrace(s.GetContext(), errors.New("sort: Payload is not array"))
	}

	slices.SortStableFunc[any](slice, func(i, j any) bool {
		if err != nil {
			return false
		}

		payload := data.NewPayload()

		for k, v := range s.GetAll() {
			payload.Set(k, v, true)
		}

		payload.Set(s.Config.ActionSetter.Key, map[string]interface{}{"a": i, "b": j}, s.Config.ActionSetter.Override)

		err = s.Next(payload)

		if err != nil {
			return false
		}

		res := payload.Get(s.Config.Result.ActionGetter.Key)
		ok, returnValue := res.(bool)
		if !ok {
			return false
		}

		return returnValue
	})
	return
}

func (s *Sort) GetConfig() interface{} {
	return s.Config
}

func (s *Sort) SetPayload(data tailsafe.DataInterface) {
	s.DataInterface = data
}

func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(Sort)
	p.StepInterface = step
	p.Config = new(Config)
	return p
}
