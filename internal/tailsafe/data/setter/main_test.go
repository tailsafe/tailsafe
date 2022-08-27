package setter

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var (
	example = `
{
}`
)

func TestValidator(t *testing.T) {
	data := map[string]any{}

	err := Set(data).Apply("", "hello")
	if !assert.NotNil(t, err) {
		return
	}

	assert.Equal(t, "key (``) cannot be null", err.Error())
}

func TestGetWithInvalidReflectData(t *testing.T) {
	err := Set(nil).Apply("message", "hello")

	assert.NotNil(t, err)
}

func TestSetSimple(t *testing.T) {
	data := map[string]any{}

	err := Set(data).Apply("message", "hello")
	if !assert.Nil(t, err) {
		return
	}

	assert.Equal(t, "hello", data["message"])
	t.Logf("value=%s", data)

}

func TestSimpleMap(t *testing.T) {
	data := map[string]any{}
	err := json.Unmarshal([]byte(example), &data)

	if !assert.Nil(t, err) {
		return
	}
	err = Set(data).Apply("message.say", "hello")
	if !assert.Nil(t, err) {
		return
	}

	assert.Equal(t, map[string]any{"message": map[string]any{"say": "hello"}}, data)
	t.Logf("value=%s", data)
}

func TestSetSliceWithOverride(t *testing.T) {
	data := []any{map[string]any{"hello": "world"}}

	err := Set(data).
		SetOverride(true).
		Apply("[0].hello", "me")
	if !assert.Nil(t, err) {
		return
	}

	assert.Equal(t, data, []interface{}{map[string]interface{}{"hello": "me"}})
}

func TestSetSliceWithoutOverride(t *testing.T) {
	data := []any{map[string]any{"hello": "world"}}

	err := Set(data).
		SetOverride(false).
		Apply("[0].hello", "me")

	if !assert.NotNil(t, err) {
		return
	}

	assert.Equal(t, err.Error(), "impossible to set the value `me` because the key `[0].hello` already contains this value `world`")
}

func TestCreateSetSlice(t *testing.T) {
	data := []any{map[string]any{}}

	err := Set(data).
		SetOverride(false).
		Apply("[].hello", "me")

	if !assert.Nil(t, err) {
		return
	}

	log.Print(data)

	assert.Equal(t, data, []interface{}{map[string]interface{}{"hello": "me"}})
}

func TestCreateSetSliceWithInvalidIndex(t *testing.T) {
	var data []any

	err := Set(data).
		SetOverride(false).
		Apply("[1].hello", "me")

	if !assert.NotNil(t, err) {
		return
	}

	assert.Equal(t, "index 1, len(0) not exist, please use [] for append", err.Error())
}
func TestCreateSetSliceAppend(t *testing.T) {
	data := []any{map[string]any{}}
	data = append(data, map[string]any{"hello": "world"})

	err := Set(data).
		SetOverride(false).
		Apply("[].hello", "me")

	if !assert.Nil(t, err) {
		return
	}

	log.Print(data)

	assert.Equal(t, []interface{}{map[string]interface{}{"hello": "me"}, map[string]interface{}{"hello": "world"}}, data)
}
