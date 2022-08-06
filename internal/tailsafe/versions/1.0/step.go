// Package __0 supports versions 1.0 template files.
package __0

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tailsafe/tailsafe/internal/tailsafe/actions"
	"github.com/tailsafe/tailsafe/internal/tailsafe/modules"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"github.com/tidwall/gjson"
	"os"
	"regexp"
	"strings"
	"sync"
)

// Step represents a single action.
// Title is the name of the action.
// Use indicates the name of the action to be used to perform the process
type Step struct {
	context.Context
	sync.Mutex

	Engine   tailsafe.EngineInterface
	Title    string   `yaml:"title"`
	Use      string   `yaml:"use"`
	LogLevel string   `yaml:"log-level"`
	Config   any      `yaml:"config"`
	Key      string   `yaml:"key"`
	Needs    []string `yaml:"needs"`

	StepsGeneric []tailsafe.StepInterface
	Steps        []*Step `yaml:"steps"`

	StepsSuccessGeneric []tailsafe.StepInterface
	StepsSuccess        []*Step `yaml:"if-action-success"`

	StepsFailGeneric []tailsafe.StepInterface
	StepsFail        []*Step `yaml:"if-action-fail"`

	Async   bool     `yaml:"async"`
	Wait    []string `yaml:"wait"`
	Payload tailsafe.DataInterface

	NextIsAlreadyCalled bool
}

func (s *Step) SetTitle(title string) tailsafe.StepInterface {
	s.Title = title
	return s
}

func (s *Step) SetConfig(data any) tailsafe.StepInterface {
	s.Config = data
	return s
}

func (s *Step) GetChildLevel() int {
	return s.Engine.GetChildLevel()
}

func (s *Step) GetEngine() tailsafe.EngineInterface {
	return s.Engine
}

func (s *Step) IsAsync() bool {
	return s.Async
}

func (s *Step) GetWait() []string {
	return s.Wait
}

// GetLogLevel returns the log level for the step.
func (s *Step) GetLogLevel() int {
	switch strings.ToLower(s.LogLevel) {
	case "info":
		return tailsafe.LOG_INFO
	case "verbose":
		return tailsafe.LOG_VERBOSE
	case "none":
		return tailsafe.LOG_NONE
	}
	return tailsafe.LOG_INFO
}

// Resolve resolves value with path into
func (s *Step) Resolve(path string, data map[string]any) any {
	b, err := json.Marshal(data)
	if err != nil {
		return path
	}
	var re = regexp.MustCompile(`(?m){{(.*)}}`)
	result := re.FindAllStringSubmatch(path, -1)
	if len(result) == 0 {
		return path
	}
	value := gjson.Get(string(b), strings.TrimSpace(result[0][1]))
	if value.Type == gjson.String {
		return strings.ReplaceAll(path, result[0][0], fmt.Sprintf("%v", value.String()))
	}
	return value.Value()
}
func (s *Step) SetEngine(engine tailsafe.EngineInterface) {
	s.Engine = engine
}
func (s *Step) SetContext(ctx context.Context) {
	s.Context = ctx
}
func (s *Step) GetContext() context.Context {
	return s.Context
}
func (s *Step) Next(payload tailsafe.DataInterface) tailsafe.ErrActionInterface {
	// if action use next if-success or if-fail
	s.NextIsAlreadyCalled = true

	return s._Next(payload, nil, true)
}
func (s *Step) _Next(payload tailsafe.DataInterface, err error, fromAction bool) tailsafe.ErrActionInterface {
	if !fromAction && s.NextIsAlreadyCalled {
		if err == nil {
			return nil
		}
		return err.(tailsafe.ErrActionInterface)
	}
	s.Engine.EntrySubStage()
	defer func() {
		s.Engine.ExitSubStage()
	}()
	var steps []tailsafe.StepInterface

	if err != nil && len(s.GetFailSteps()) == 0 {
		return err.(tailsafe.ErrActionInterface)
	}

	if err != nil && len(s.GetFailSteps()) != 0 {
		steps = s.GetFailSteps()
	}

	if err == nil && len(s.GetSuccessSteps()) != 0 {
		steps = s.GetSuccessSteps()
	}

	for _, sub := range steps {

		sub.SetContext(s.Context)
		sub.SetPayload(payload)
		sub.SetEngine(s.Engine)

		err := sub.Call()
		if err != nil {
			return tailsafe.CatchStackTrace(s.GetContext(), err)
		}
	}
	return nil
}

func (s *Step) Plugin() tailsafe.ActionInterface {
	return nil
}

func (s *Step) SetPayload(data tailsafe.DataInterface) {
	s.Payload = data
}

func (s *Step) GetPayload() tailsafe.DataInterface {
	return s.Payload
}

// GetSteps returns the steps for this step.
func (s *Step) GetSteps() []tailsafe.StepInterface {
	if len(s.StepsGeneric) == 0 {
		for _, test := range s.Steps {
			s.StepsGeneric = append(s.StepsGeneric, test)
		}
	}
	return s.StepsGeneric
}

func (s *Step) GetSuccessSteps() []tailsafe.StepInterface {
	if len(s.StepsSuccessGeneric) == 0 {
		for _, step := range s.StepsSuccess {
			s.StepsSuccessGeneric = append(s.StepsSuccessGeneric, step)
		}
	}
	return s.StepsSuccessGeneric
}

func (s *Step) GetFailSteps() []tailsafe.StepInterface {
	if len(s.StepsFailGeneric) == 0 {
		for _, step := range s.StepsFail {
			s.StepsFailGeneric = append(s.StepsFailGeneric, step)
		}
	}
	return s.StepsFailGeneric
}

// GetTitle returns the title for this step.
func (s *Step) GetTitle() string {
	return s.Title
}

// GetUse returns the use for this step.
func (s *Step) GetUse() string {
	return strings.ReplaceAll(s.Use, "~", os.Getenv("HOME"))
}

func (s *Step) SetUse(use string) tailsafe.StepInterface {
	s.Use = use
	return s
}

// GetKey returns the key for this step.
func (s *Step) GetKey() string {
	var re = regexp.MustCompile(`(?m)%(.*)%`)
	res := re.FindAllStringSubmatch(s.Key, -1)
	for _, v := range res {
		s.Key = strings.ReplaceAll(s.Key, v[0], fmt.Sprintf("%v", s.Payload.Get(v[1])))
	}
	return strings.TrimSpace(s.Key)
}
func (s *Step) GetConfig() interface{} {
	return s.Config
}

func (s *Step) Begin() tailsafe.StageMonitoringInterface {
	// increment the stage
	s.Engine.NewStage()
	// return the stage monitor
	return NewStepMonitoring(s.Engine.GetStage())
}

func (s *Step) Call() (err error) {
	// Begin monitoring the stage
	stageMonitoring := s.Begin()
	stageLevel := s.Engine.GetChildLevel()

	defer func() {
		if s.IsAsync() {
			return
		}
		stageMonitoring.End()
		if err != nil {
			return
		}
		modules.Get[tailsafe.EventsInterface]("Events").Trigger(tailsafe.NewActionExitEvent(s, stageMonitoring, stageLevel))
	}()

	actionFunc, err := actions.Get(s.GetUse())
	if err != nil {
		return err
	}

	modules.Get[tailsafe.EventsInterface]("Events").Trigger(tailsafe.NewActionBeforeConfigureStepEvent(s))

	if s.Payload == nil {
		s.Payload = tailsafe.NewPayload()
	}

	// only set the key if it is not empty
	if s.GetKey() != "" {
		modules.Get[tailsafe.EventsInterface]("Events").Trigger(tailsafe.NewActionHasStoringKeyEvent(s))
	}

	if len(s.GetWait()) > 0 {
		modules.Get[tailsafe.EventsInterface]("Events").Trigger(tailsafe.NewActionHasWaitEvent(s))
		err = modules.Get[tailsafe.AsyncQueue]("AsyncQueue").WaitActions(s.GetWait()...)
		if err != nil {
			return
		}
		// reset monitoring for real started action
		stageMonitoring.Reset()
	}

	actionInstance := actionFunc(s)
	action := actionInstance

	out, err := json.Marshal(s.GetConfig())
	if err != nil {
		return err
	}

	tmp := action.GetConfig()

	err = json.Unmarshal(out, &tmp)
	if err != nil {
		return
	}

	// Inject need into current payload
	need := s.Engine.ExtractGlobal(s.Needs)
	for key, value := range need {
		s.GetPayload().Set(key, value)
	}

	// inject into action
	action.SetPayload(s.GetPayload())

	// configure the action
	err = action.Configure()
	if err != nil {
		return
	}
	modules.Get[tailsafe.EventsInterface]("Events").Trigger(tailsafe.NewActionAfterConfigureStepEvent(s, action))

	// if data mocked !
	// @todo check if not already set in loop mode
	mock := s.Engine.GetMockDataByKey(s.GetKey())
	if mock != nil {
		s.Engine.SetData(s.GetKey(), mock)
		return
	}

	if s.IsAsync() {
		modules.Get[tailsafe.EventsInterface]("Events").Trigger(tailsafe.NewActionIsAsyncEvent(s))
	}

	payload := func() (err error) {
		// execute the action
		err = action.Execute()
		if err != nil {
			return
		}

		// if store key is defined, store the value
		key := s.GetKey()
		result := action.GetResult()
		if strings.TrimSpace(key) != "" && result != nil {
			reservedKey := []string{"args"}
			if strings.Contains(strings.Join(reservedKey, " "), key) {
				return fmt.Errorf("key `%s` is reserved from the system [%s]", key, strings.Join(reservedKey, ", "))
			}

			// Set global data
			s.Engine.SetData(key, result)

			// Set payload state data
			s.GetPayload().Set(key, result)

			modules.Get[tailsafe.EventsInterface]("Events").Trigger(tailsafe.NewActionHasStoringDataEvent(s, result))
		}

		if !s.IsAsync() {
			return
		}

		stageMonitoring.End()
		if err != nil {
			modules.Get[tailsafe.EventsInterface]("Events").Trigger(tailsafe.NewActionExitWithErrorEvent(s, err, stageMonitoring, stageLevel+1))
			return
		}
		modules.Get[tailsafe.EventsInterface]("Events").Trigger(tailsafe.NewActionExitEvent(s, stageMonitoring, stageLevel+1))
		return
	}

	if s.IsAsync() {
		modules.Get[tailsafe.AsyncQueue]("AsyncQueue").AddActionToQueue(s.GetKey(), payload)
		return
	}

	return s._Next(s.GetPayload(), payload(), false)
}
