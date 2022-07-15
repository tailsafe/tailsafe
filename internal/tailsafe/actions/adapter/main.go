package adapteraction

import (
	"errors"
	"fmt"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"golang.org/x/exp/slices"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type AdapterAction struct {
	tailsafe.StepInterface
	Config any

	global map[string]interface{}
	data   any
	sync.Mutex
}

type ObjectType struct {
	Type       string
	Value      any
	Extra      any
	Properties map[string]interface{}
}

// New creates a new AdapterAction
func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(AdapterAction)
	p.StepInterface = step
	return p
}

/* Generic Getters */

// GetData returns the data for the action
func (ta *AdapterAction) GetData() interface{} {
	return ta.data
}

// GetConfig returns the configuration for the action
func (ta *AdapterAction) GetConfig() interface{} {
	if ta.Config == nil {
		return ta.Config
	}
	return ta.Config
}

/* Generics setters */

// SetConfig sets the configuration for the action
func (ta *AdapterAction) SetConfig(config any) {
	if config == nil {
		return
	}
	ta.Config = config
}

// SetGlobal sets the global data for the action
func (ta *AdapterAction) SetGlobal(data map[string]interface{}) {
	ta.global = data
}

// Configure the action
func (ta *AdapterAction) Configure() tailsafe.ErrActionInterface {
	return nil
}

// Execute the action
func (ta *AdapterAction) Execute() (err tailsafe.ErrActionInterface) {
	return ta.parse(ta.GetConfig())
}

func (ta *AdapterAction) SetInternalGlobal(key string, data any) {
	ta.Lock()
	defer ta.Unlock()

	ta.global[key] = data
	return
}

func (ta *AdapterAction) GetInternalGlobal(key string) any {
	ta.Lock()
	defer ta.Unlock()

	return ta.global[key]
}

func (ta *AdapterAction) parse(data any) tailsafe.ErrActionInterface {
	obj, err := ta.getType(data)
	if err != nil {
		return tailsafe.CatchStackTrace(ta.GetContext(), err)
	}
	switch obj.Type {
	case "array":
		arrayData, err := ta.parseArray(obj)
		if err != nil {
			return tailsafe.CatchStackTrace(ta.GetContext(), err)
		}
		ta.data = arrayData
	case "object":
		ta.data = make(map[string]interface{})
		return ta.parseProperties(obj.Properties, ta.data.(map[string]any))
	}
	return nil
}
func (ta *AdapterAction) parseArray(obj *ObjectType) (data []any, err tailsafe.ErrActionInterface) {
	d := ta.Resolve(fmt.Sprintf("%v", obj.Value), ta.global)
	if d == nil {
		return nil, tailsafe.CatchStackTrace(ta.GetContext(), fmt.Errorf("could not resolve %v", obj.Value))
	}
	rf := reflect.TypeOf(d).Kind()
	if rf != reflect.Slice {
		return data, tailsafe.CatchStackTrace(ta.GetContext(), errors.New(rf.String()+" is not a slice"))
	}
	tmp := ta.GetInternalGlobal("this")
	for _, this := range d.([]interface{}) {
		ta.SetInternalGlobal("this", this)
		newObject := make(map[string]interface{})
		err := ta.parseProperties(obj.Properties, newObject)
		if err != nil {
			return data, tailsafe.CatchStackTrace(ta.GetContext(), err)
		}
		data = append(data, newObject)
	}
	ta.SetInternalGlobal("this", tmp)
	return
}
func (ta *AdapterAction) parseProperties(properties map[string]interface{}, dst map[string]interface{}) tailsafe.ErrActionInterface {
	for k, v := range properties {
		obj, err := ta.getType(v)
		if err != nil {
			return tailsafe.CatchStackTrace(ta.GetContext(), err)
		}
		switch obj.Type {
		case "array":
			data, err := ta.parseArray(obj)
			if err != nil {
				return tailsafe.CatchStackTrace(ta.GetContext(), err)
			}
			ta.Lock()
			dst[k] = data
			ta.Unlock()
		case "object":
			newObject := make(map[string]interface{})
			err := ta.parseProperties(obj.Properties, newObject)
			if err != nil {
				return tailsafe.CatchStackTrace(ta.GetContext(), err)
			}
			ta.Lock()
			dst[k] = newObject
			ta.Unlock()
		case "string":
			ta.Lock()
			dst[k] = fmt.Sprintf("%v", ta.Resolve(fmt.Sprintf("%v", obj.Value), ta.global))
			ta.Unlock()
		case "number":
			value := ta.Resolve(fmt.Sprintf("%v", obj.Value), ta.global)
			n, err := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 64)
			if err != nil {
				return tailsafe.CatchStackTrace(ta.GetContext(), err)
			}
			ta.Lock()
			dst[k] = n
			ta.Unlock()
		case "boolean":
			value := ta.Resolve(fmt.Sprintf("%v", obj.Value), ta.global)
			b, err := strconv.ParseBool(fmt.Sprintf("%v", value))
			if err != nil {
				return tailsafe.CatchStackTrace(ta.GetContext(), err)
			}
			ta.Lock()
			dst[k] = b
			ta.Unlock()
		case "datetime":
			value := fmt.Sprintf("%+v", ta.Resolve(fmt.Sprintf("%v", obj.Value), ta.global))
			layout := time.RFC3339
			timeT, err := time.Parse(layout, value)
			if err != nil {
				return tailsafe.CatchStackTrace(ta.GetContext(), err)
			}
			ta.Lock()
			dst[k] = timeT.Format(time.RFC3339)
			ta.Unlock()
		}
	}
	return nil
}

// getType returns the primary type of the data
func (ta *AdapterAction) getType(data any) (objType *ObjectType, err error) {
	m, ok := data.(map[string]interface{})
	if !ok {
		err = errors.New("config is not an object")
		return
	}
	objType = new(ObjectType)
	objType.Type, ok = m["type"].(string)
	if !ok {
		err = fmt.Errorf("type must be in string, not %v (%v)", reflect.TypeOf(m["type"]), m["type"])
		return
	}
	objType.Properties, ok = m["properties"].(map[string]interface{})
	if !ok {
		objType.Properties = make(map[string]interface{})
	}
	objType.Value, _ = m["value"]
	objType.Extra, _ = m["extra"]

	authorized := []string{"string", "number", "object", "array", "datetime", "boolean"}

	if !slices.Contains(authorized, objType.Type) {
		err = fmt.Errorf("type `%s` is not supported", objType.Type)
		return
	}

	return
}
