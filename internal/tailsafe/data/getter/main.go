package getter

import (
	"reflect"
	"strconv"
	"strings"
)

type Getter struct {
	key         string
	focusedList []string
	index       *int
}

func Get(focused string, data any) any {

	if strings.TrimSpace(focused) == "" {
		return nil
	}

	e := new(Getter)
	e.key = focused

	e._Analyse()

	r := reflect.ValueOf(data)

	switch r.Kind() {
	case reflect.Map:
		return e._ExtractObjectValue(r)
	case reflect.Slice:
		return e._ExtractArrayValue(data.([]any))
	default:
		return nil
	}
}
func (g *Getter) _Parse(s string) bool {
	if s[0] != '[' {
		return false
	}
	s = s[1:]

	if s[len(s)-1] != ']' {
		return false
	}

	s = s[:len(s)-1]

	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}

	g.index = new(int)
	*g.index, _ = strconv.Atoi(s)

	return true
}
func (g *Getter) _Analyse() {
	g.focusedList = strings.Split(strings.TrimSpace(g.key), ".")
	g.key = g.focusedList[0]

	g._Parse(g.key)
}
func (g *Getter) _ExtractObjectValue(data reflect.Value) any {
	focused := data.MapIndex(reflect.ValueOf(g.key))
	if focused.IsValid() {
		if len(g.focusedList) > 1 {
			return Get(strings.Join(g.focusedList[1:], "."), focused.Interface())
		}
		return focused.Interface()
	}
	return nil
}
func (g *Getter) _ExtractArrayValue(data []any) any {
	if g.index == nil {
		return nil
	}
	if len(data) < *g.index {
		return nil
	}
	if len(g.focusedList) == 1 {
		return data[*g.index]
	}
	return Get(strings.Join(g.focusedList[1:], "."), data[*g.index])
}
