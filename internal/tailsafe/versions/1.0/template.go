package __0

import (
	"fmt"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"gopkg.in/yaml.v3"
	"strings"
)

type Template struct {
	Version     string `yaml:"versions"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Maintainer  string `yaml:"maintainer"`
	Revision    int    `yaml:"revision"`

	// stdout is the output of the template
	StdOut []string `yaml:"stdout"`

	// Flags
	Args []struct {
		Type     string `json:"type"`
		Name     string `json:"name"`
		Usage    string `json:"usage"`
		Default  any    `json:"default"`
		Required bool   `json:"required"`
		Value    any
	} `json:"args"`

	// Steps are the steps in the template
	StepInterface []tailsafe.StepInterface
	Steps         []*Step `yaml:"steps"`
}

// GetVersion returns the versions of the template
func (t Template) GetVersion() string {
	return t.Version
}

// GetTitle returns the title of the template
func (t Template) GetTitle() string {
	return t.Title
}

// GetDescription returns the description of the template
func (t Template) GetDescription() string {
	return t.Description
}

// GetRevision returns the revision of the template
func (t Template) GetRevision() int {
	return t.Revision
}

func (t Template) GetMaintainer() string {
	return t.Maintainer
}

// GetSteps returns the steps in the template
func (t *Template) GetSteps() []tailsafe.StepInterface {
	return t.StepInterface
}

// GetStdOut returns the stdout of the template
func (t Template) GetStdOut() []string {
	return t.StdOut
}

func (t Template) GetDependencies() []string {
	return t.getSteps(t.GetSteps())
}

// recursive function to get use of all steps
func (t *Template) getSteps(steps []tailsafe.StepInterface) []string {
	var allSteps []string
	for _, step := range steps {
		allSteps = append(allSteps, step.GetUse())
		allSteps = append(allSteps, t.getSteps(step.GetSteps())...)
	}
	return allSteps
}

func (t *Template) InjectPreStep(steps []tailsafe.StepInterface) {
	for _, step := range steps {
		t.StepInterface = append([]tailsafe.StepInterface{step.(tailsafe.StepInterface)}, t.StepInterface...)
	}
}

func (t *Template) InjectPostStep(_ []tailsafe.StepInterface) {
}

func (t *Template) NewStep() tailsafe.StepInterface {
	return new(Step)
}

// SetEnv configure the environment for the template
func (t Template) SetEnv(args string) (data any, err error) {
	split := strings.Split(args, ",")
	keyValue := make(map[string]any)

	for _, arg := range split {
		splitArg := strings.Split(strings.TrimSpace(arg), ":")
		if len(splitArg) != 2 {
			continue
		}

		keyValue[splitArg[0]] = splitArg[1]
	}

	if len(keyValue) == 0 {
		return
	}

	for k, f := range t.Args {
		switch f.Type {
		case "string":
			t.Args[k].Value = keyValue[f.Name]
		}
	}

	// check if all required flags are set
	var argsRequired []string
	var finalArgs = make(map[string]any)
	for _, f := range t.Args {
		if !f.Required {
			continue
		}
		if f.Value != nil {
			finalArgs[f.Name] = f.Value
			continue
		}
		argsRequired = append(argsRequired, fmt.Sprintf("%s (%s)", f.Name, f.Usage))
	}

	// trigger error if required flags are not set
	if len(argsRequired) > 0 {
		err = fmt.Errorf("required flags are not set: %s", strings.Join(argsRequired, ", "))
		return
	}

	return finalArgs, nil
}

// Parse parses the template from the given yaml
func Parse(data []byte) (template any, err error) {
	tmp := new(Template)
	err = yaml.Unmarshal(data, &tmp)
	if err != nil {
		return
	}

	for _, test := range tmp.Steps {
		tmp.StepInterface = append(tmp.StepInterface, test)
	}
	// force template type
	template = tmp
	return
}
