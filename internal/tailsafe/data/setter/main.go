package setter

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Setter struct {
	key      string
	origin   []string
	value    any
	override bool

	focusedList []string
	index       *int
	append      bool
	data        any
}

func Set(data any) (se *Setter) {
	se = new(Setter)
	se.data = data

	return
}

// SetOrigin muse be user to keep key origin
func (se *Setter) SetOrigin(origin []string) *Setter {
	se.origin = origin

	return se
}

// SetValue must be used to set the value
func (se *Setter) SetValue(value any) *Setter {
	se.value = value

	return se
}

// SetKey must be used to set the key
func (se *Setter) SetKey(key string) *Setter {
	se.key = key

	return se
}

func (se *Setter) SetOverride(value bool) *Setter {
	se.override = value

	return se
}

// Apply try to set value
func (se *Setter) Apply(key string, value any) (err error) {
	se.key = key
	se.value = value

	if strings.TrimSpace(se.key) == "" {
		return fmt.Errorf("key (`%v`) cannot be empty", se.key)
	}

	data := reflect.ValueOf(se.data)

	if !data.IsValid() {
		err = fmt.Errorf("data is null")
		return
	}

	se._Analyse()

	switch data.Kind() {
	case reflect.Map:
		se.data, err = se._SetObjectKeyValue(data)
	case reflect.Slice:
		se.data, err = se._SetSliceValue(data)
	}

	return
}

func (se *Setter) _Parse(s string) bool {
	if s == "[]" {
		se.append = true
	}

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

	se.index = new(int)
	*se.index, _ = strconv.Atoi(s)

	return true
}

func (se *Setter) _Analyse() {
	se.focusedList = strings.Split(strings.TrimSpace(se.key), ".")
	se.key = se.focusedList[0]

	se._Parse(se.key)
}

func (se *Setter) _SetObjectKeyValue(data reflect.Value) (res any, err error) {

	if se.append {
		err = fmt.Errorf("append key detected (`%v`), data found has not a slice (`%v`)", se.key, data.Type())
		return
	}

	focused := data.MapIndex(reflect.ValueOf(se.key))
	if focused.IsValid() && !se.override {
		return focused.Interface(), fmt.Errorf("impossible to set the value `%v` because the key `%v` already contains this value `%s`", se.value, strings.Join(append(se.origin, []string{se.key}...), "."), focused.Interface())
	}

	var nData any
	if focused.IsValid() {
		nData = focused.Interface()
	} else {
		nData = make(map[string]any)
	}

	if len(se.focusedList) > 1 {

		err = Set(nData).
			SetOverride(se.override).
			SetOrigin(append(se.origin, []string{se.key}...)).
			Apply(strings.Join(se.focusedList[1:], "."), se.value)

		if err != nil {
			return
		}

		data.SetMapIndex(reflect.ValueOf(se.key), reflect.ValueOf(nData))

		res = data.Interface()

		return
	}

	data.SetMapIndex(reflect.ValueOf(se.key), reflect.ValueOf(se.value))
	res = data.Interface()

	return
}

func (se *Setter) _SetSliceValue(data reflect.Value) (res any, err error) {
	if se.index == nil {
		return
	}
	if *se.index > data.Len()-1 {
		err = fmt.Errorf("index %v, len(%d) not exist, please use [] for append", *se.index, data.Len())
		return
	}

	value := data.Index(*se.index)
	if !value.IsValid() {
		return
	}

	err = Set(value.Interface()).
		SetOverride(se.override).
		SetOrigin(append(se.origin, []string{se.key}...)).
		Apply(strings.Join(se.focusedList[1:], "."), se.value)

	if err != nil {
		return
	}

	res = data.Interface()

	return
}
