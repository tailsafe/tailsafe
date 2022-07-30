package adapteraction

import (
	"github.com/stretchr/testify/assert"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"github.com/tailsafe/tailsafe/pkg/tailsafe/setuphelper"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	step := setuphelper.NewStepTesting()
	action := New(step)

	assert.NotNil(t, action)
}

func TestAdapterAction_Configure(t *testing.T) {
	step := setuphelper.NewStepTesting()
	action := New(step)

	assert.NotNil(t, action)

	t.Run("Empty config", func(t *testing.T) {
		err := action.Configure()
		assert.Nil(t, err)
	})
}

func TestAdapterAction_Execute(t *testing.T) {
	step := setuphelper.NewStepTesting()
	action := New(step)

	assert.NotNil(t, action)

	t.Run("should work with object adapter", func(t *testing.T) {
		var data []byte
		data, err := os.ReadFile("./testdata/simple_begin_by_object.yml")
		if !assert.Nil(t, err) {
			return
		}

		cfg := action.GetConfig()
		err = yaml.Unmarshal(data, cfg)

		if !assert.Nil(t, err) {
			return
		}

		payload := tailsafe.NewPayload()
		payload.Set("{{ test.number }}", 0)
		payload.Set("{{ test.string }}", "test")
		payload.Set("{{ test.datetime }}", "2020-01-01T00:00:00Z")
		payload.Set("{{ test.boolean }}", true)

		action.SetPayload(payload)

		err = action.Configure()
		if !assert.Nil(t, err) {
			return
		}

		err = action.Execute()
		if !assert.Nil(t, err) {
			return
		}

		assert.Equal(t, true, action.GetResult().(map[string]interface{})["boolean"])
		assert.Equal(t, "2020-01-01T00:00:00Z", action.GetResult().(map[string]interface{})["datetime"])
		assert.Equal(t, int64(0), action.GetResult().(map[string]interface{})["number"])
		assert.Equal(t, "test", action.GetResult().(map[string]interface{})["string"])

		assert.Equal(t, true, action.GetResult().(map[string]interface{})["sub_object"].(map[string]interface{})["boolean"])
		assert.Equal(t, "2020-01-01T00:00:00Z", action.GetResult().(map[string]interface{})["sub_object"].(map[string]interface{})["datetime"])
		assert.Equal(t, int64(0), action.GetResult().(map[string]interface{})["sub_object"].(map[string]interface{})["number"])
		assert.Equal(t, "test", action.GetResult().(map[string]interface{})["sub_object"].(map[string]interface{})["string"])
	})
}
