package versions

import (
	"errors"
	__0 "github.com/tailsafe/tailsafe/internal/tailsafe/versions/1.0"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
)

// TemplateInterface contains the interfaces needed to generate the code
type TemplateInterface interface {
	/* Getters */

	GetTitle() string
	GetDescription() string
	GetMaintainer() string
	GetRevision() int

	GetSteps() []tailsafe.StepInterface
	GetDependencies() []string

	GetStdOut() []string

	NewStep() tailsafe.StepInterface

	InjectPreStep([]tailsafe.StepInterface)
	InjectPostStep([]tailsafe.StepInterface)

	/* Setters */

	// SetEnv sets the environment for the template
	SetEnv(env string) (data any, err error)
}

// List of all available versions
var templates = map[string]func(data []byte) (any, error){
	"1.0": __0.Parse,
}

// GetTemplate returns the template for the given versions.
func GetTemplate(data []byte, label string) (templateInterface TemplateInterface, err error) {
	parse, ok := templates[label]

	if !ok {
		err = errors.New("No template found for versions " + label)
		return
	}

	v, err := parse(data)
	if err != nil {
		return
	}

	templateInterface = v.(TemplateInterface)

	return
}
