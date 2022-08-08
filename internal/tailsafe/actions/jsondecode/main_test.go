package jsonDecodeAction

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"github.com/tailsafe/tailsafe/pkg/tailsafe/setuphelper"
	"testing"
)

func TestNew(t *testing.T) {
	step := setuphelper.NewStepTesting()
	action := New(step)

	assert.NotNil(t, action)
}

func TestJsonDecodeAction_Configure(t *testing.T) {
	step := setuphelper.NewStepTesting()
	action := New(step)

	assert.NotNil(t, action)

	t.Run("Empty config", func(t *testing.T) {
		err := action.Configure()
		assert.NotNil(t, err)

		assert.Contains(t, err.Error(), "JsonDecodeAction: Value cannot be empty")
	})
	t.Run("Valid config", func(t *testing.T) {
		cfg := action.GetConfig()
		assert.NotNil(t, cfg)

		cfg.(*Config).Value = `{"value":"test"}`

		err := action.Configure()
		assert.Nil(t, err)
	})
}

func TestJsonEncodeAction_Execute(t *testing.T) {
	step := setuphelper.NewStepTesting()
	action := New(step)

	assert.NotNil(t, action)

	t.Run("With string value", func(t *testing.T) {
		cfg := action.GetConfig()
		assert.NotNil(t, cfg)

		cfg.(*Config).Value = `{"value":"test"}`

		action.SetPayload(tailsafe.NewPayload())

		err := action.Configure()
		assert.Nil(t, err)

		err = action.Execute()
		assert.Nil(t, err)

		assert.Equal(t, map[string]interface{}{"value": "test"}, action.GetResult())
	})

	t.Run("With resolved string json value", func(t *testing.T) {
		cfg := action.GetConfig()
		assert.NotNil(t, cfg)

		cfg.(*Config).Value = `value?`

		payload := tailsafe.NewPayload()
		payload.Set("value", `{"value":"test"}`)

		action.SetPayload(payload)

		err := action.Configure()
		assert.Nil(t, err)

		err = action.Execute()
		assert.Nil(t, err)

		assert.Equal(t, map[string]interface{}{"value": "test"}, action.GetResult())
	})
	t.Run("With resolved invalid type", func(t *testing.T) {
		cfg := action.GetConfig()
		assert.NotNil(t, cfg)

		cfg.(*Config).Value = `{{ value }}`

		payload := tailsafe.NewPayload()
		payload.Set("{{ value }}", map[string]interface{}{"value": "test"})

		action.SetPayload(payload)

		err := action.Configure()
		assert.Nil(t, err)

		err = action.Execute()
		assert.NotNil(t, err)

		assert.IsType(t, err, tailsafe.ErrAction{})
	})

	t.Run("With resolved invalid string value", func(t *testing.T) {
		cfg := action.GetConfig()
		assert.NotNil(t, cfg)

		cfg.(*Config).Value = `{{ value }}`

		payload := tailsafe.NewPayload()
		payload.Set("{{ value }}", "test")

		action.SetPayload(payload)

		err := action.Configure()
		assert.Nil(t, err)

		err = action.Execute()
		assert.NotNil(t, err)

		assert.IsType(t, err, tailsafe.ErrAction{})
		assert.IsType(t, err.GetOriginal(), &json.SyntaxError{})
	})
}
