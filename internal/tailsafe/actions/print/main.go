package printAction

import (
	"errors"
	"fmt"
	"github.com/tailsafe/tailsafe/internal/tailsafe/modules"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"strings"
)

type Config struct {
	Format string `yaml:"format"`
	Values []any  `yaml:"values"`
}

type PrintAction struct {
	tailsafe.DataInterface
	tailsafe.StepInterface
	Config *Config

	global tailsafe.DataInterface
	result string
}

func (pa *PrintAction) Configure() (err tailsafe.ErrActionInterface) {
	if strings.TrimSpace(pa.Config.Format) == "" {
		return tailsafe.CatchStackTrace(pa.GetContext(), errors.New("PrintAction: Format cannot be empty"))
	}
	return
}

func (pa *PrintAction) Execute() (err tailsafe.ErrActionInterface) {
	for k, value := range pa.Config.Values {
		v, ok := value.(string)
		if !ok {
			continue
		}
		pa.Config.Values[k] = pa.Resolve(v, pa.GetAll())
	}

	modules.GetEventsModule().Trigger(tailsafe.NewActionStdoutEvent(pa.StepInterface, fmt.Sprintf(pa.Config.Format, pa.Config.Values...), pa.GetChildLevel()))
	return
}
func (pa *PrintAction) GetResult() interface{} {
	return pa.result
}
func (pa *PrintAction) GetConfig() interface{} {
	return pa.Config
}
func (pa *PrintAction) SetPayload(data tailsafe.DataInterface) {
	pa.DataInterface = data
}
func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(PrintAction)
	p.StepInterface = step
	p.Config = new(Config)
	return p
}
