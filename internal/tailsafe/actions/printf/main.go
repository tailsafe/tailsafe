package printfAction

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

type PrintfAction struct {
	tailsafe.DataInterface
	tailsafe.StepInterface
	Config *Config

	global tailsafe.DataInterface
	result string
}

func (pa *PrintfAction) Configure() (err tailsafe.ErrActionInterface) {
	if strings.TrimSpace(pa.Config.Format) == "" {
		return tailsafe.CatchStackTrace(pa.GetContext(), errors.New("PrintfAction: Format cannot be empty"))
	}
	return
}

func (pa *PrintfAction) Execute() (err tailsafe.ErrActionInterface) {
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
func (pa *PrintfAction) GetResult() interface{} {
	return pa.result
}
func (pa *PrintfAction) GetConfig() interface{} {
	return pa.Config
}
func (pa *PrintfAction) SetPayload(data tailsafe.DataInterface) {
	pa.DataInterface = data
}
func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(PrintfAction)
	p.StepInterface = step
	p.Config = new(Config)
	return p
}
