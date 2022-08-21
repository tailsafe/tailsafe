package tailsafe

import "context"

// StepInterface is the interface that must be implemented by all steps.
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
	GetNextSteps() []StepInterface
	GetSuccessSteps() []StepInterface
	GetFailSteps() []StepInterface

	HasFailed() bool

	GetEngine() EngineInterface
	GetContext() context.Context
	GetPayload() DataInterface
	Next(payload DataInterface) ErrActionInterface
	Resolve(path string, data map[string]any) any
	IsAsync() bool

	/* Public Setters */

	SetContext(context.Context)
	SetEngine(EngineInterface)
	SetPayload(DataInterface)

	SetUse(string) StepInterface
	SetTitle(string) StepInterface
	SetConfig(data any) StepInterface
}
