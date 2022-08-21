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
	printAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/print"
	"github.com/tailsafe/tailsafe/internal/tailsafe/actions/replace"
	setterAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/setter"
	sortAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/sort"
	templateAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/template"
	termsAction "github.com/tailsafe/tailsafe/internal/tailsafe/actions/terms"
	"github.com/tailsafe/tailsafe/internal/tailsafe/modules"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"log"
	"plugin"
	"strings"
	"sync"
)

type Actions struct {
	sync.Mutex
	data map[string]func(runtime tailsafe.StepInterface) tailsafe.ActionInterface
}

// create singleton with a factory method
var instance *Actions

// init is called when the package is loaded.
func init() {
	// create singleton
	instance = &Actions{
		data: make(map[string]func(runtime tailsafe.StepInterface) tailsafe.ActionInterface),
	}

	// Only internals packages of Golang for the internal actions.
	instance.Lock()
	defer instance.Unlock()

	instance.data["internal/setter"] = setterAction.New
	instance.data["internal/map"] = mapAction.New
	instance.data["internal/foreach"] = foreachAction.New
	instance.data["internal/terms"] = termsAction.New
	instance.data["internal/datetime"] = datetimeAction.New
	instance.data["internal/replace"] = replaceAction.New
	instance.data["internal/http"] = httpAction.New
	instance.data["internal/template"] = templateAction.New
	instance.data["internal/adapter"] = adapterAction.New
	instance.data["internal/exec"] = execAction.New
	instance.data["internal/payload"] = payloadAction.New
	instance.data["internal/print"] = printAction.New
	instance.data["internal/jsonEncode"] = jsonEncodeAction.New
	instance.data["internal/jsonDecode"] = jsonDecodeAction.New
	instance.data["internal/sort"] = sortAction.New
}
func Get(name string) (action func(runtime tailsafe.StepInterface) tailsafe.ActionInterface, err error) {
	// if dev action
	if strings.HasPrefix(name, "/") && instance.data[name] == nil {
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
		instance.data[name] = action
		instance.Unlock()
	}
	return _GetInternal(name)
}

func _GetInternal(name string) (action func(runtime tailsafe.StepInterface) tailsafe.ActionInterface, err error) {
	instance.Lock()
	defer instance.Unlock()

	var ok bool
	action, ok = instance.data[name]
	if !ok {
		err = &tailsafe.ErrActionNotFound{Name: name}
		return
	}
	return
}
