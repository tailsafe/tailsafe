package setuphelper

import (
	"context"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
)

type Step struct {
	context.Context
	tailsafe.DataInterface
	title    string
	config   any
	use      string
	ley      string
	key      string
	logLevel int
	wait     []string
	steps    []Step
	payload  map[string]any
}

func NewStepTesting() tailsafe.StepInterface {
	s := new(Step)
	s.payload = make(map[string]any)
	return s
}

func (s *Step) Call() (err error) {
	return
}

func (s *Step) GetTitle() string {
	return s.title
}

func (s *Step) GetKey() string {
	return s.key
}

func (s *Step) GetUse() string {
	return s.use
}

func (s *Step) GetLogLevel() int {
	return s.logLevel
}

func (s *Step) GetChildLevel() int {
	return 0
}

func (s *Step) GetWait() []string {
	return s.wait
}

func (s *Step) GetSteps() []tailsafe.StepInterface {
	return []tailsafe.StepInterface{}
}

func (s *Step) GetSuccessSteps() []tailsafe.StepInterface {
	return []tailsafe.StepInterface{}
}

func (s *Step) GetFailSteps() []tailsafe.StepInterface {
	return []tailsafe.StepInterface{}
}

func (s *Step) GetEngine() tailsafe.EngineInterface {
	return nil
}

func (s *Step) GetContext() context.Context {
	return s.Context
}

func (s *Step) GetPayload() tailsafe.DataInterface {
	return s.DataInterface
}

func (s *Step) Next(payload tailsafe.DataInterface) tailsafe.ErrActionInterface {
	return nil
}

func (s *Step) Resolve(path string, data map[string]any) any {
	v, ok := data[path]
	if !ok {
		return path
	}
	return v
}

func (s *Step) IsAsync() bool {
	return false
}

func (s *Step) SetContext(ctx context.Context) {
	s.Context = ctx
}

func (s *Step) SetEngine(engineInterface tailsafe.EngineInterface) {
}

func (s *Step) SetPayload(dataInterface tailsafe.DataInterface) {
}

func (s *Step) SetUse(use string) tailsafe.StepInterface {
	s.use = use

	return s
}

func (s *Step) SetTitle(s2 string) tailsafe.StepInterface {
	s.title = s2

	return s
}

func (s *Step) SetConfig(data any) tailsafe.StepInterface {
	s.config = data

	return s
}

func (s *Step) SetResolve(key string, value any) tailsafe.StepInterface {
	s.payload[key] = value

	return s
}
