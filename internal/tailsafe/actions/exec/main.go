package execaction

import (
	"bufio"
	"errors"
	"github.com/tailsafe/tailsafe/internal/tailsafe/modules"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"os"
	"os/exec"
	"strings"
)

type Config struct {
	Command string `yaml:"command"`
	Path    string `yaml:"path"`
}
type ExecAction struct {
	tailsafe.StepInterface
	Config *Config

	global       map[string]interface{}
	data         string
	commandSlice []string
}

func (ex *ExecAction) Configure() (err tailsafe.ErrActionInterface) {
	if ex.Config == nil {
		return tailsafe.CatchStackTrace(ex.GetContext(), errors.New("ExecAction: Config is nil"))
	}

	ex.Config.Command = strings.TrimSpace(ex.Config.Command)
	if ex.Config.Command == "" {
		return tailsafe.CatchStackTrace(ex.GetContext(), errors.New("ExecAction: Command is empty"))
	}

	ex.commandSlice = strings.Split(ex.Config.Command, " ")
	return
}

func (ex *ExecAction) Execute() tailsafe.ErrActionInterface {
	cmd := exec.CommandContext(ex.GetContext(), ex.commandSlice[0], ex.commandSlice[1:]...)
	cmd.Env = os.Environ()

	if ex.Config.Path != "" {
		cmd.Dir = ex.Config.Path
	}

	stdout, err := cmd.StdoutPipe()
	rd := bufio.NewReader(stdout)

	go func() {
		indentLevel := ex.GetChildLevel()
		for {
			str, err := rd.ReadString('\n')
			if err != nil {
				break
			}
			modules.GetEventsModule().Trigger(tailsafe.NewActionStdoutEvent(ex.StepInterface, str, indentLevel))
		}
	}()

	err = cmd.Run()

	if err != nil {
		return tailsafe.CatchStackTrace(ex.GetContext(), err)
	}
	return nil
}
func (ex *ExecAction) GetData() interface{} {
	return ex.data
}
func (ex *ExecAction) GetConfig() interface{} {
	if ex.Config == nil {
		return &Config{}
	}
	return ex.Config
}
func (ex *ExecAction) SetConfig(config interface{}) {
	ex.Config = config.(*Config)
}
func (ex *ExecAction) SetGlobal(data map[string]interface{}) {
	ex.global = data
}
func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(ExecAction)
	p.StepInterface = step
	return p
}
