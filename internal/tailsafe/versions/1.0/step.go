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
	"regexp"
	"strings"
)

// Step represents a single action.
// Title is the name of the action.
// Use indicates the name of the action to be used to perform the process
type Step struct {
	context.Context
	Engine   tailsafe.EngineInterface
	Title    string                   `yaml:"title"`
	Use      string                   `yaml:"use"`
	LogLevel string                   `yaml:"log-level"`
	Data     any                      `yaml:"data"`
	Key      string                   `yaml:"key"`
	Needs    []string                 `yaml:"needs"`
	Steps    []tailsafe.StepInterface `yaml:"steps"`
	Async    bool                     `yaml:"async"`
	Wait     []string                 `yaml:"wait"`
	Current  any
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
func (s *Step) Resolve(path string, data any) any {
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
func (s *Step) Next() tailsafe.ErrActionInterface {
	s.Engine.EntrySubStage()
	defer s.Engine.ExitSubStage()

	for _, sub := range s.GetSteps() {
		sub.SetContext(s.Context)
		sub.SetCurrent(s.Current)
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

func (s *Step) SetCurrent(data interface{}) {
	s.Current = data
}

// GetSteps returns the steps for this step.
func (s Step) GetSteps() []tailsafe.StepInterface {
	return s.Steps
}

// GetTitle returns the title for this step.
func (s Step) GetTitle() string {
	return s.Title
}

// GetUse returns the use for this step.
func (s Step) GetUse() string {
	return s.Use
}

// GetKey returns the key for this step.
func (s Step) GetKey() string {
	var re = regexp.MustCompile(`(?m)%(.*)%`)
	res := re.FindAllStringSubmatch(s.Key, -1)
	for _, v := range res {
		data, ok := s.Current.(map[string]interface{})
		if ok {
			s.Key = strings.ReplaceAll(s.Key, v[0], fmt.Sprintf("%v", data[v[1]]))
		}
	}
	return strings.TrimSpace(s.Key)
}
func (s Step) GetData() interface{} {
	return s.Data
}

func (s Step) Begin() tailsafe.StageMonitoringInterface {
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

	out, err := json.Marshal(s.GetData())
	if err != nil {
		return err
	}

	tmp := action.GetConfig()
	err = json.Unmarshal(out, &tmp)
	if err != nil {
		return
	}

	// Set the config into the action
	action.SetConfig(tmp)

	// extract global with need!
	extractGlobal := s.Engine.ExtractGlobal(s.Needs)

	// set current with each context
	extractGlobal["current"] = s.Current

	// inject into action
	action.SetGlobal(extractGlobal)

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
		//s.Engine.Log(tailsafe.NAMESPACE_WORKFLOW, tailsafe.LOG_VERBOSE, " ↳ Uses mock data : %v", s.Engine.Pretty(mock))
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
		if s.GetKey() != "" {
			data := action.GetData()
			// if not null, saving !
			if data != nil {
				// secure key name
				reservedKey := []string{"args"}
				if strings.Contains(strings.Join(reservedKey, " "), s.GetKey()) {
					return fmt.Errorf("key `%s` is reserved from the system [%s]", s.GetKey(), strings.Join(reservedKey, ", "))
				}

				s.Engine.SetData(s.GetKey(), data)
				modules.Get[tailsafe.EventsInterface]("Events").Trigger(tailsafe.NewActionHasStoringDataEvent(s, data))
			}
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

	return payload()
}
