package resolver

import (
	"strconv"
	"strings"
)

type resolver struct {
	focused     string
	focusedList []string
	key         *int64
}

func Get(focused string, data any) any {

	if strings.TrimSpace(focused) == "" {
		return nil
	}

	e := new(resolver)
	e.focused = focused

	e._Analyse()

	// Dispatch search
	switch d := data.(type) {
	case map[string]interface{}:
		return e._ExtractObjectValue(d)
	case []interface{}:
		return e._ExtractArrayValue(d)
	default:
		return nil
	}
}
func (e *resolver) _IsInt(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
func (e *resolver) _Analyse() {
	e.focusedList = strings.Split(strings.TrimSpace(e.focused), ".")
	e.focused = e.focusedList[0]

	if e._IsInt(e.focused) {
		e.key = new(int64)
		*e.key, _ = strconv.ParseInt(e.focused, 10, 64)
		return
	}
}
func (e *resolver) _ExtractObjectValue(data map[string]interface{}) interface{} {
	if val, ok := data[e.focused]; ok {
		if len(e.focusedList) > 1 {
			return Get(strings.Join(e.focusedList[1:], "."), val)
		}
		return val
	}
	return nil
}
func (e *resolver) _ExtractArrayValue(data []any) any {
	if e.key == nil {
		return nil
	}
	if len(data) < int(*e.key) {
		return nil
	}
	if len(e.focusedList) == 1 {
		return data[*e.key]
	}
	return Get(strings.Join(e.focusedList[1:], "."), data[*e.key])
}
