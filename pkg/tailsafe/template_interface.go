package tailsafe

// TemplateInterface is the interface that must be implemented by all template actions.
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
