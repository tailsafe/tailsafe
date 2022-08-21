package setterAction

import (
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"reflect"
)

type Config struct {
	ActionGetter tailsafe.ActionGetter `json:"action-getters"`
	ActionSetter tailsafe.ActionSetter `json:"action-setters"`
}

type Setter struct {
	tailsafe.StepInterface
	tailsafe.DataInterface

	Config *Config
}

func (s *Setter) Configure() (err tailsafe.ErrActionInterface) {
	gv := s.Resolve(s.Config.ActionGetter.Key, s.GetAll())

	value := s.Resolve(s.Config.ActionSetter.Value, s.GetAll())
	tf := reflect.ValueOf(value)

	if tf.Kind() == reflect.Map {
		tf.SetMapIndex(reflect.ValueOf(s.Config.ActionSetter.Key), reflect.ValueOf(gv))
	}
	return
}

func (s *Setter) GetResult() any {
	return s.Config
}

func (s *Setter) Execute() (err tailsafe.ErrActionInterface) {
	return
}
func (s *Setter) GetConfig() any {
	return &s.Config
}
func (s *Setter) SetPayload(data tailsafe.DataInterface) {
	s.DataInterface = data
}
func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(Setter)
	p.StepInterface = step
	return p
}
