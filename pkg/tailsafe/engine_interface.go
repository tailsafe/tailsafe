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

type Utils interface {
	Pretty(v any, level int) any
}

// EngineInterface represents a engine in the workflow
type EngineInterface interface {
	DataInterface

	ExtractGlobal(required []string) map[string]any

	SetPath(path string) EngineInterface
	SetEnv(env string) EngineInterface
	SetDataPath(path string) EngineInterface

	Run()
	/* Mock Payload process */

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

// Error returns the error message with more Payload
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
