package actions

import (
	"fmt"
	"github.com/tailsafe/tailsafe/internal/tailsafe/actions/If"
	adapterAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/adapter"
	"github.com/tailsafe/tailsafe/internal/tailsafe/actions/datetime"
	execAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/exec"
	httpAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/http"
	"github.com/tailsafe/tailsafe/internal/tailsafe/actions/loop"
	"github.com/tailsafe/tailsafe/internal/tailsafe/actions/replace"
	templateAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/template"
	"github.com/tailsafe/tailsafe/internal/tailsafe/modules"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"log"
	"os"
	"plugin"
	"strings"
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
	// if internal action
	if strings.HasPrefix(name, "internal/") {
		return _GetInternal(name)
	}
	// if dev action
	if strings.HasPrefix(name, "/") {
		split := strings.Split(name, "/")
		name := fmt.Sprintf("%s@dev.so", split[len(split)-1])

		log.Print("[DEV] Build Action: ", name)

		log.Print(modules.GetUtilsModule().GetAppActionDir())

		pp, err := plugin.Open(fmt.Sprintf("%s/%s", modules.GetUtilsModule().GetAppActionDir(), name))
		if err != nil {
			log.Fatalln(err)
		}

		pl, err := pp.Lookup("New")
		if err != nil {
			log.Fatalln(err)
		}

		c := pl.(func(runtime tailsafe.StepInterface) tailsafe.ActionInterface)

		log.Print(c)
	}

	os.Exit(1)
	/*	var ok bool
		action, ok = instance.data[name]
		if !ok {
			err = &tailsafe.ErrActionNotFound{Name: name}
			return
		}*/
	return
}

func _GetInternal(name string) (action func(runtime tailsafe.StepInterface) tailsafe.ActionInterface, err error) {
	var ok bool
	action, ok = instance.data[name]
	if !ok {
		err = &tailsafe.ErrActionNotFound{Name: name}
		return
	}
	return
}
