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
	tailsafe.DataInterface

	Config any

	data any
	sync.Mutex
}

type ObjectType struct {
	Type       string
	Value      any
	Extra      any
	Nullable   bool
	Properties map[string]any
	Items      map[string]any
}

// New creates a new AdapterAction
func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(AdapterAction)
	p.StepInterface = step

	return p
}

/* Generic Getters */

// GetResult returns the data for the action
func (ta *AdapterAction) GetResult() any {
	return ta.data
}

// GetConfig returns the configuration for the action
func (ta *AdapterAction) GetConfig() any {
	return &ta.Config
}

// SetPayload sets the global data for the action
func (ta *AdapterAction) SetPayload(data tailsafe.DataInterface) {
	ta.DataInterface = data
}

// Configure the action
func (ta *AdapterAction) Configure() tailsafe.ErrActionInterface {
	return nil
}

// Execute the action
func (ta *AdapterAction) Execute() (err tailsafe.ErrActionInterface) {
	return ta.parse(ta.Config)
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
		ta.data = make(map[string]any)
		return ta.parseProperties(obj.Properties, ta.data.(map[string]any))
	}
	return nil
}
func (ta *AdapterAction) parseArray(obj *ObjectType) (data []any, err tailsafe.ErrActionInterface) {
	d := ta.Resolve(fmt.Sprintf("%v", obj.Value), ta.GetAll())
	if d == nil {
		return nil, tailsafe.CatchStackTrace(ta.GetContext(), fmt.Errorf("could not resolve %v", obj.Value))
	}
	rf := reflect.TypeOf(d).Kind()
	if rf != reflect.Slice {
		return data, tailsafe.CatchStackTrace(ta.GetContext(), fmt.Errorf("%v (%s) is not a slice", d, rf.String()))
	}
	tmp := ta.Get(tailsafe.THIS)

	for _, this := range d.([]any) {

		ta.Set(tailsafe.THIS, this)

		if len(obj.Items) > 0 {
			newObject := make(map[string]any)
			err := ta.parseProperties(map[string]any{tailsafe.THIS: obj.Items}, newObject)
			if err != nil {
				return data, tailsafe.CatchStackTrace(ta.GetContext(), err)
			}
			data = append(data, newObject[tailsafe.THIS])
			continue
		}

		newObject := make(map[string]any)
		err := ta.parseProperties(obj.Properties, newObject)
		if err != nil {
			return data, tailsafe.CatchStackTrace(ta.GetContext(), err)
		}

		data = append(data, newObject)
	}

	ta.Set(tailsafe.THIS, tmp)
	return
}
func (ta *AdapterAction) parseProperties(properties map[string]any, dst map[string]any) tailsafe.ErrActionInterface {
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
			newObject := make(map[string]any)
			err := ta.parseProperties(obj.Properties, newObject)
			if err != nil {
				return tailsafe.CatchStackTrace(ta.GetContext(), err)
			}
			ta.Lock()
			dst[k] = newObject
			ta.Unlock()
		case "string":
			ta.Lock()
			value := ta.Resolve(fmt.Sprintf("%v", obj.Value), ta.GetAll())
			if value == nil && !obj.Nullable {
				return tailsafe.CatchStackTrace(ta.GetContext(), fmt.Errorf("could not resolve %v", obj.Value))
			}
			dst[k] = fmt.Sprintf("%v", value)
			ta.Unlock()
		case "number":
			value := ta.Resolve(fmt.Sprintf("%v", obj.Value), ta.GetAll())
			if value == nil && !obj.Nullable {
				return tailsafe.CatchStackTrace(ta.GetContext(), fmt.Errorf("could not resolve %v", obj.Value))
			}
			n, err := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 64)
			if err != nil {
				return tailsafe.CatchStackTrace(ta.GetContext(), err)
			}
			ta.Lock()
			dst[k] = n
			ta.Unlock()
		case "boolean":
			value := ta.Resolve(fmt.Sprintf("%v", obj.Value), ta.GetAll())
			if value == nil && !obj.Nullable {
				return tailsafe.CatchStackTrace(ta.GetContext(), fmt.Errorf("could not resolve %v", obj.Value))
			}
			b, err := strconv.ParseBool(fmt.Sprintf("%v", value))
			if err != nil {
				return tailsafe.CatchStackTrace(ta.GetContext(), err)
			}
			ta.Lock()
			dst[k] = b
			ta.Unlock()
		case "datetime":
			value := fmt.Sprintf("%+v", ta.Resolve(fmt.Sprintf("%v", obj.Value), ta.GetAll()))
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
	m, ok := data.(map[string]any)
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
	objType.Properties, ok = m["properties"].(map[string]any)
	if !ok {
		objType.Properties = make(map[string]any)
	}
	objType.Value, _ = m["value"]
	objType.Extra, _ = m["extra"]

	authorized := []string{"string", "number", "datetime", "boolean", "object", "array"}

	if !slices.Contains(authorized, objType.Type) {
		err = fmt.Errorf("type `%s` is not supported", objType.Type)
		return
	}

	objType.Items, ok = m["items"].(map[string]any)
	if !ok {
		objType.Items = make(map[string]any)
	}

	return
}
