package textReader

import (
	"bufio"
	"github.com/tailsafe/tailsafe/internal/tailsafe/data"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"os"
)

type Config struct {
	Path         string                `json:"path"`
	ActionSetter tailsafe.ActionSetter `json:"action-setter"`
}

type TextReader struct {
	tailsafe.DataInterface
	tailsafe.StepInterface
	Config *Config

	result any
}

func (fa *TextReader) Configure() (err tailsafe.ErrActionInterface) {
	return
}

func (fa *TextReader) Execute() (err tailsafe.ErrActionInterface) {
	path := fa.Resolve(fa.Config.Path, fa.GetAll())
	file, osErr := os.Open(path.(string))

	if osErr != nil {
		return tailsafe.CatchStackTrace(fa.GetContext(), osErr)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		payload := data.NewPayload()
		payload.Set(fa.Config.ActionSetter.Key, scanner.Text(), fa.Config.ActionSetter.Override)

		err = fa.Next(payload)
		if err != nil {
			return err
		}
	}

	file.Close()
	return
}

func (fa *TextReader) GetResult() interface{} {
	return fa.result
}

func (fa *TextReader) GetConfig() interface{} {
	return fa.Config
}

func (fa *TextReader) SetPayload(data tailsafe.DataInterface) {
	fa.DataInterface = data
}

func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(TextReader)
	p.StepInterface = step
	p.Config = new(Config)
	return p
}
