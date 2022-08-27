package stringAction

import (
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"strings"
)

type Config struct {
	Type         string                `json:"type"`
	ActionGetter tailsafe.ActionGetter `json:"action-getter"`
	ActionSetter tailsafe.ActionSetter `json:"action-setter"`
	Extra        string                `json:"extra"`
}

type String struct {
	tailsafe.DataInterface
	tailsafe.StepInterface
	Config []Config
}

func (s *String) Configure() (err tailsafe.ErrActionInterface) {
	return
}

func (s *String) Execute() (err tailsafe.ErrActionInterface) {
	for _, c := range s.Config {
		switch strings.ToLower(c.Type) {
		case "trim":
			s.Set(c.ActionSetter.Key, strings.ToLower(s.GetValue(c.ActionGetter.Value)), c.ActionSetter.Override)
		case "lowercase":
			s.Set(c.ActionSetter.Key, strings.ToLower(s.GetValue(c.ActionGetter.Value)), c.ActionSetter.Override)
		case "split":
			s.Set(c.ActionSetter.Key, strings.Split(s.GetValue(c.ActionGetter.Value), c.Extra), c.ActionSetter.Override)
		}
	}
	return
}

func (s *String) GetValue(resolve string) string {
	value := s.Resolve(resolve, s.GetAll())
	v, ok := value.(string)
	if !ok {
		return ""
	}
	return v
}

func (s *String) GetResult() interface{} {
	return nil
}

func (s *String) GetConfig() interface{} {
	return &s.Config
}

func (s *String) SetPayload(data tailsafe.DataInterface) {
	s.DataInterface = data
}

func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(String)
	p.StepInterface = step
	p.Config = []Config{}
	return p
}
