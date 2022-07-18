package tailsafe

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

const (
	LOG_NONE    = iota
	LOG_INFO    = 1
	LOG_VERBOSE = 2

	NAMESPACE_DEFAULT  = "default"
	NAMESPACE_WORKFLOW = "global"
)

type Logger interface {
	Log(LoggerPayload)
}
type LoggerPayload interface {
	GetNamespace() string
	GetLevel() int
	GetMessage() string
	GetArgs() []any

	SetNamespace(string) LoggerPayload
	SetLevel(int) LoggerPayload
	SetMessage(string) LoggerPayload
	SetArgs(...any) LoggerPayload
}

type TemplateInterface interface {
	NewStep() StepInterface

	GetTitle() string
	GetDescription() string
	GetMaintainer() string
	GetRevision() int

	GetSteps() []StepInterface
	GetDependencies() []string

	InjectPreStep([]StepInterface)
	InjectPostStep([]StepInterface)
}

// StepInterface represents a step in the workflow
type StepInterface interface {
	/* Public Getters */

	Call() (err error)
	GetTitle() string
	GetKey() string
	GetUse() string
	GetLogLevel() int
	GetChildLevel() int
	GetWait() []string
	GetSteps() []StepInterface
	GetEngine() EngineInterface
	GetContext() context.Context
	Next() (err ErrActionInterface)
	Resolve(path string, data any) any
	IsAsync() bool

	/* Public Setters */

	SetContext(context.Context)
	SetEngine(EngineInterface)
	SetCurrent(data any)

	SetUse(string) StepInterface
	SetTitle(string) StepInterface
	SetData(data any) StepInterface
}

// ActionInterface represents an action in the workflow
type ActionInterface interface {
	// Configure the action
	Configure() (err ErrActionInterface)

	// Execute the action
	Execute() (err ErrActionInterface)

	/* Setters */

	// SetConfig sets the configuration for the action
	SetConfig(any)
	SetGlobal(data map[string]interface{})

	/* Getters */

	// GetConfig returns the configuration for the action
	GetConfig() any
	GetData() any
}

type EventsInterface interface {
	Trigger(event EventInterface)
}
type Utils interface {
	Pretty(v any, level int) any
}
type Event interface {
	SetBody(...any)
	SetKey(string)

	GetBody() []any
	GetKey() string
}

// EngineInterface represents a engine in the workflow
type EngineInterface interface {
	ExtractGlobal(required []string) map[string]any
	SetData(key string, data any)

	SetPath(path string) EngineInterface
	SetEnv(env string) EngineInterface
	SetPathData(path string) EngineInterface

	Run()

	/* Mock data process */

	GetMockDataByKey(key string) any

	// NewStage creates a new stage
	NewStage()
	// GetStage returns the current stage
	GetStage() int

	GetChildLevel() int

	// EntrySubStage enters a sub stage
	EntrySubStage()
	// ExitSubStage exits a sub stage
	ExitSubStage()
}

/* Async process */

// AsyncQueue is a generic func for running an async action.
type AsyncQueue interface {

	// AddActionToQueue Add an action to the asynchronous queue
	AddActionToQueue(action string, call func() error) AsyncQueue

	// WaitActions Allows you to wait for the end of the specified actions
	WaitActions(actions ...string) error

	// WaitAll Allows you to wait for the end of all actions
	WaitAll()

	// SetWorkers sets the number of workers for the asynchronous queue
	SetWorkers(num int) AsyncQueue
}

type StageMonitoringInterface interface {
	GetStage() int
	GetStageDuration() time.Duration
	Reset()
	End()
}

/*
	Public error types for the tailsafe-cli package.
*/

// ErrActionNotFound is returned when an action is not found.
type ErrActionNotFound struct {
	Name string
}

func (e *ErrActionNotFound) Error() string {
	return fmt.Sprintf("Action `%s` cannot be found", e.Name)
}

// ErrContinue is returned when a step should continue.
type ErrContinue struct {
	Message string
}

func (e *ErrContinue) SetError(err error) {
	e.Message = err.Error()
}

func (e *ErrContinue) Error() string {
	return e.Message
}

type ErrActionInterface interface {
	error
	GetStackTrace() string
	GetOriginal() error
}

type ErrAction struct {
	trace    string
	caller   string
	line     int
	file     string
	original error
}

// GetOriginal returns the original error without the stack trace.
func (e ErrAction) GetOriginal() error {
	return e.original
}

// GetStackTrace returns the stack trace of the error.
func (e ErrAction) GetStackTrace() string {
	return e.trace
}

// Error returns the error message with more data
func (e ErrAction) Error() string {
	return fmt.Sprintf("`%s` was triggered by %s:%d from %s", e.original.Error(), e.caller, e.line, e.file)
}

// CatchStackTrace is generic func for returning an error.
func CatchStackTrace(_ context.Context, err error) ErrAction {
	var Err ErrAction
	if errors.As(err, &Err) {
		return err.(ErrAction)
	}

	est := ErrAction{}
	est.original = err

	sp := strings.Split(string(debug.Stack()), "\n")
	est.trace = strings.TrimSpace(strings.Join(sp[6:], "\n"))
	var line int
	var file string
	pc, file, line, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok {
		est.caller = details.Name()
		est.line = line
		est.file = file
	}

	return est
}
