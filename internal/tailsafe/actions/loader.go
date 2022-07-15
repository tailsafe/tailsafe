package actions

import (
	"github.com/tailsafe/tailsafe/internal/tailsafe/actions/If"
	adapterAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/adapter"
	"github.com/tailsafe/tailsafe/internal/tailsafe/actions/datetime"
	execAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/exec"
	httpAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/http"
	"github.com/tailsafe/tailsafe/internal/tailsafe/actions/loop"
	"github.com/tailsafe/tailsafe/internal/tailsafe/actions/replace"
	templateAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/template"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
)

type Actions struct {
	data map[string]func(runtime tailsafe.StepInterface) tailsafe.ActionInterface
}

// create singleton with a factory method
var instance *Actions

// init is called when the package is loaded.
func init() {
	instance = &Actions{
		data: make(map[string]func(runtime tailsafe.StepInterface) tailsafe.ActionInterface),
	}

	// Only internals packages of Golang for the internal actions.
	// @todo: add a way to add autoload of external actions.

	instance.data["internal/for"] = loop.New
	instance.data["internal/if"] = If.New
	instance.data["internal/datetime"] = datetime.New
	instance.data["internal/replace"] = replace.New
	instance.data["internal/http"] = httpAction.New
	instance.data["internal/template"] = templateAction.New
	instance.data["internal/adapter"] = adapterAction.New
	instance.data["internal/exec"] = execAction.New
}
func Get(name string) (action func(runtime tailsafe.StepInterface) tailsafe.ActionInterface, err error) {
	var ok bool
	action, ok = instance.data[name]
	if !ok {
		err = &tailsafe.ErrActionNotFound{Name: name}
		return
	}
	return
}
