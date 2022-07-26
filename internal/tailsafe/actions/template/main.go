package templateaction

import (
	"bytes"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"text/template"
)

type Config struct {
	Template string `yaml:"template"`
}

type TemplateAction struct {
	tailsafe.StepInterface
	Config *Config

	global tailsafe.DataInterface
	data   string
}

func (ta *TemplateAction) Configure() (err tailsafe.ErrActionInterface) {
	return
}

func (ta *TemplateAction) Execute() (err tailsafe.ErrActionInterface) {
	t, herr := template.New(ta.GetTitle()).Parse(ta.Config.Template)
	if herr != nil {
		return tailsafe.CatchStackTrace(ta.GetContext(), herr)
	}
	var tpl bytes.Buffer
	herr = t.Execute(&tpl, ta.global)
	if herr != nil {
		return tailsafe.CatchStackTrace(ta.GetContext(), herr)
	}

	ta.data = tpl.String()
	return
}
func (ta *TemplateAction) GetResult() interface{} {
	return ta.data
}
func (ta *TemplateAction) GetConfig() interface{} {
	if ta.Config == nil {
		return &Config{}
	}
	return ta.Config
}
func (ta *TemplateAction) SetPayload(data tailsafe.DataInterface) {
	ta.global = data
}
func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(TemplateAction)
	p.StepInterface = step
	return p
}
