package actions

import (
	"fmt"
	adapterAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/adapter"
	datetimeAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/datetime"
	execAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/exec"
	foreachAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/foreach"
	httpAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/http"
	jsonDecodeAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/jsondecode"
	jsonEncodeAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/jsonencode"
	mapAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/map"
	payloadAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/payload"
	printfAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/printf"
	setAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/set"
	sortAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/sort"
	stringAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/string"
	templateAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/template"
	termsAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/terms"
	textreader "github.com/tailsafe/tailsafe/internal/tailsafe/actions/textreader"
	"github.com/tailsafe/tailsafe/internal/tailsafe/modules"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"log"
	"plugin"
	"strings"
	"sync"
)

type ActionLoader struct {
	sync.Mutex
	actions map[string]func(runtime tailsafe.StepInterface) tailsafe.ActionInterface
}

// create singleton with a factory method
var instance *ActionLoader

// init is called when the package is loaded.
func init() {
	// create singleton
	instance = &ActionLoader{
		actions: make(map[string]func(runtime tailsafe.StepInterface) tailsafe.ActionInterface),
	}

	// Only internals packages of Golang for the internal actions.
	instance.Lock()
	defer instance.Unlock()

	instance.actions["internal/textReader"] = textreader.New
	instance.actions["internal/set"] = setAction.New
	instance.actions["internal/map"] = mapAction.New
	instance.actions["internal/foreach"] = foreachAction.New
	instance.actions["internal/terms"] = termsAction.New
	instance.actions["internal/datetime"] = datetimeAction.New
	instance.actions["internal/http"] = httpAction.New
	instance.actions["internal/template"] = templateAction.New
	instance.actions["internal/adapter"] = adapterAction.New
	instance.actions["internal/exec"] = execAction.New
	instance.actions["internal/payload"] = payloadAction.New
	instance.actions["internal/printf"] = printfAction.New
	instance.actions["internal/jsonEncode"] = jsonEncodeAction.New
	instance.actions["internal/jsonDecode"] = jsonDecodeAction.New
	instance.actions["internal/sort"] = sortAction.New
	instance.actions["internal/string"] = stringAction.New

	// always in lowercase
	for k, v := range instance.actions {
		instance.actions[strings.ToLower(k)] = v
	}
}

// Get returns the action by name.
func Get(name string) (action func(runtime tailsafe.StepInterface) tailsafe.ActionInterface, err error) {
	// more secure for search !
	name = strings.ToLower(name)
	// if dev actions
	if strings.HasPrefix(name, "/") && instance.actions[name] == nil {
		split := strings.Split(name, "/")
		pluginName := fmt.Sprintf("%s@dev.so", split[len(split)-1])

		var pp *plugin.Plugin
		pp, err = plugin.Open(fmt.Sprintf("%s/%s", modules.GetUtilsModule().GetAppActionDir(), pluginName))
		if err != nil {
			log.Fatalln(err)
		}

		var pl plugin.Symbol
		pl, err = pp.Lookup("New")
		if err != nil {
			log.Fatalln(err)
		}

		action = pl.(func(runtime tailsafe.StepInterface) tailsafe.ActionInterface)

		instance.Lock()
		instance.actions[name] = action
		instance.Unlock()
	}
	return _GetInternal(name)
}

// _GetInternal returns the action by name.
func _GetInternal(name string) (action func(runtime tailsafe.StepInterface) tailsafe.ActionInterface, err error) {
	instance.Lock()
	defer instance.Unlock()

	var ok bool
	action, ok = instance.actions[name]
	if !ok {
		err = &tailsafe.ErrActionNotFound{Name: name}
		return
	}
	return
}
