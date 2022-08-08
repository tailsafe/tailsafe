package jsonEncodeAction

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"github.com/tailsafe/tailsafe/pkg/tailsafe/setuphelper"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	step := setuphelper.NewStepTesting()
	action := New(step)

	assert.NotNil(t, action)
}

func TestJsonEncodeAction_Configure(t *testing.T) {
	step := setuphelper.NewStepTesting()
	action := New(step)

	assert.NotNil(t, action)

	t.Run("Empty config", func(t *testing.T) {
		err := action.Configure()
		assert.NotNil(t, err)

		assert.Contains(t, err.Error(), "JsonEncodeAction: Value cannot be empty")
	})

	t.Run("Valid config", func(t *testing.T) {
		config := make(map[string]interface{})
		config["value"] = "test"

		cfg := action.GetConfig()
		assert.NotNil(t, cfg)

		cfg.(*Config).Value = config["value"]

		err := action.Configure()
		assert.Nil(t, err)
	})
}

func TestJsonEncodeAction_SetPayload(t *testing.T) {
	step := setuphelper.NewStepTesting()
	action := New(step)

	assert.NotNil(t, action)

	t.Run("SetPayload", func(t *testing.T) {
		action.SetPayload(tailsafe.NewPayload())
	})
}

func TestJsonEncodeAction_Execute(t *testing.T) {
	step := setuphelper.NewStepTesting()
	action := New(step)

	assert.NotNil(t, action)

	t.Run("With object value", func(t *testing.T) {
		config := make(map[string]interface{})
		config["value"] = "test"

		cfg := action.GetConfig()
		assert.NotNil(t, cfg)

		cfg.(*Config).Value = config

		action.SetPayload(tailsafe.NewPayload())

		err := action.Configure()
		assert.Nil(t, err)

		err = action.Execute()
		assert.Nil(t, err)

		assert.Equal(t, `{"value":"test"}`, action.GetResult())
	})

	t.Run("With resolved string value", func(t *testing.T) {

		cfg := action.GetConfig()
		assert.NotNil(t, cfg)

		cfg.(*Config).Value = "value?"

		config := make(map[string]interface{})
		config["value"] = "test"
		payload := tailsafe.NewPayload()
		payload.Set("value", config)

		action.SetPayload(payload)

		err := action.Configure()
		assert.Nil(t, err)

		err = action.Execute()
		assert.Nil(t, err)

		assert.Equal(t, `{"value":"test"}`, action.GetResult())
	})

	t.Run("Should thrown a error with invalid value", func(t *testing.T) {
		cfg := action.GetConfig()
		assert.NotNil(t, cfg)

		cfg.(*Config).Value = make(chan int)

		payload := tailsafe.NewPayload()

		action.SetPayload(payload)

		err := action.Configure()
		assert.Nil(t, err)

		err = action.Execute()
		assert.NotNil(t, err)

		assert.IsType(t, err.GetOriginal(), &json.UnsupportedTypeError{Type: reflect.TypeOf(make(chan int))})
	})
}
