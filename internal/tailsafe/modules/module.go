package modules

import (
	"errors"
	"strings"
	"sync"
)

var modules map[string]any

func init() {
	modules = make(map[string]any)
}

func Register(name string, module any) {
	var m sync.Mutex
	m.Lock()
	defer m.Unlock()

	nameLower := strings.ToLower(name)

	if _, ok := modules[nameLower]; ok {
		panic("Module " + name + " already registered")
	}

	modules[nameLower] = module
}

func Get[T any](name string) (module T) {
	var m sync.Mutex
	m.Lock()
	defer m.Unlock()

	nameLower := strings.ToLower(name)

	module, ok := modules[nameLower].(T)
	if !ok {
		panic("Module " + name + " not found, use Register() to register it and check with Requires()")
	}

	return
}

func Requires(required []string) (err error) {
	var m sync.Mutex
	m.Lock()
	defer m.Unlock()

	var noRegisteredModules []string
	for _, module := range required {
		moduleLower := strings.ToLower(module)
		if _, ok := modules[moduleLower]; ok {
			continue
		}
		noRegisteredModules = append(noRegisteredModules, module)
	}

	if len(noRegisteredModules) > 0 {
		err = errors.New("No module found for names " + strings.Join(noRegisteredModules, ", "))
		return
	}
	return
}

func Reset() {
	var m sync.Mutex
	m.Lock()
	defer m.Unlock()

	modules = make(map[string]any)
}
