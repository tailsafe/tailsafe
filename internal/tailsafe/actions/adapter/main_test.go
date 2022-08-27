package adapteraction

import (
	"github.com/stretchr/testify/assert"
	"github.com/tailsafe/tailsafe/internal/tailsafe/data"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"github.com/tailsafe/tailsafe/pkg/tailsafe/setuphelper"
	"gopkg.in/yaml.v3"
	"log"
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

func TestAdapterAction_ExecuteWithObjectConfig(t *testing.T) {
	step := setuphelper.NewStepTesting()
	action := New(step)

	if !assert.NotNil(t, action, "AdapterAction should not be nil") {
		return
	}

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

	payload := data.NewPayload()
	payload.Set("object", map[string]interface{}{
		"number":   0,
		"string":   "test",
		"datetime": "2020-01-01T00:00:00Z",
		"boolean":  "true",
	})
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

}
func TestAdapterAction_ExecuteWithObjectConfigAndEmptyPayload(t *testing.T) {
	step := setuphelper.NewStepTesting()
	action := New(step)

	if !assert.NotNil(t, action, "AdapterAction should not be nil") {
		return
	}

	var data []byte
	data, err := os.ReadFile("./testdata/simple_begin_by_object.yml")
	if !assert.NoError(t, err) {
		return
	}

	cfg := action.GetConfig()
	err = yaml.Unmarshal(data, cfg)

	if !assert.NoError(t, err) {
		return
	}

	payload := data.NewPayload()
	action.SetPayload(payload)

	err = action.Configure()
	if !assert.NoError(t, err) {
		return
	}

	err = action.Execute()
	if !assert.Error(t, err) {
		return
	}
	assert.ErrorContains(t, err, "could not resolve")
}

func TestAdapterAction_ExecuteWithArrayConfig(t *testing.T) {
	step := setuphelper.NewStepTesting()
	action := New(step)

	if !assert.NotNil(t, action, "AdapterAction should not be nil") {
		return
	}

	var data []byte
	data, err := os.ReadFile("./testdata/simple_begin_by_array.yml")
	if !assert.Nil(t, err) {
		return
	}

	cfg := action.GetConfig()
	err = yaml.Unmarshal(data, cfg)

	if !assert.Nil(t, err) {
		return
	}

	payload := data.NewPayload()
	payload.Set("global", map[string]interface{}{
		"number":   999,
		"string":   "test",
		"datetime": "2020-01-01T00:00:00Z",
		"boolean":  "true",
		"array": []map[string]interface{}{
			{
				"number":   0,
				"boolean":  "true",
				"string":   "test",
				"datetime": "2020-01-01T00:00:00Z",
			},
			{
				"number":   1,
				"boolean":  "false",
				"string":   "test-1",
				"datetime": "2020-01-02T00:00:00Z",
			},
		},
	})

	action.SetPayload(payload)

	err = action.Configure()
	if !assert.Nil(t, err) {
		return
	}

	err = action.Execute()
	if !assert.Nil(t, err) {
		return
	}

	assert.Equal(t, true, len(action.GetResult().([]interface{})) == 2)
}
func TestAdapterAction_ExecuteWithNoResolveArray(t *testing.T) {
	step := setuphelper.NewStepTesting()
	action := New(step)

	if !assert.NotNil(t, action, "AdapterAction should not be nil") {
		return
	}

	var data []byte
	data, err := os.ReadFile("./testdata/bad_type_in_properties.yml")
	if !assert.Nil(t, err) {
		return
	}

	cfg := action.GetConfig()
	err = yaml.Unmarshal(data, cfg)

	if !assert.Nil(t, err) {
		return
	}

	payload := data.NewPayload()
	action.SetPayload(payload)

	err = action.Configure()
	if !assert.NoError(t, err) {
		return
	}

	err = action.Execute()

	if !assert.NotNil(t, err) {
		return
	}

	assert.IsType(t, tailsafe.ErrAction{}, err)
	assert.ErrorContains(t, err, "could not resolve global.array?")
}

func TestAdapterAction_ExecuteArrayBadType(t *testing.T) {
	step := setuphelper.NewStepTesting()
	action := New(step)

	if !assert.NotNil(t, action, "AdapterAction should not be nil") {
		return
	}

	var data []byte
	data, err := os.ReadFile("./testdata/bad_type_in_properties.yml")
	if !assert.Nil(t, err) {
		return
	}

	cfg := action.GetConfig()
	err = yaml.Unmarshal(data, cfg)

	if !assert.Nil(t, err) {
		return
	}

	payload := data.NewPayload()
	payload.Set("global", map[string]interface{}{
		"array": "test",
	})
	action.SetPayload(payload)

	err = action.Configure()
	if !assert.NoError(t, err) {
		return
	}

	err = action.Execute()

	if !assert.NotNil(t, err) {
		return
	}
	log.Print(err)
	assert.IsType(t, tailsafe.ErrAction{}, err)
	assert.ErrorContains(t, err, "test (string) is not a slice")
}

func TestAdapterAction_ExecuteWithBadTypeForType(t *testing.T) {
	step := setuphelper.NewStepTesting()
	action := New(step)

	if !assert.NotNil(t, action, "AdapterAction should not be nil") {
		return
	}

	var data []byte
	data, err := os.ReadFile("./testdata/bad_type_in_properties.yml")
	if !assert.Nil(t, err) {
		return
	}

	cfg := action.GetConfig()
	err = yaml.Unmarshal(data, cfg)

	if !assert.Nil(t, err) {
		return
	}

	payload := data.NewPayload()
	payload.Set("global", map[string]interface{}{
		"number":   999,
		"string":   "test",
		"datetime": "2020-01-01T00:00:00Z",
		"boolean":  "true",
		"array": []map[string]interface{}{
			{
				"number":   0,
				"boolean":  "true",
				"string":   "test",
				"datetime": "2020-01-01T00:00:00Z",
			},
			{
				"number":   1,
				"boolean":  "false",
				"string":   "test-1",
				"datetime": "2020-01-02T00:00:00Z",
			},
		},
	})

	action.SetPayload(payload)

	err = action.Configure()
	if !assert.NoError(t, err) {
		return
	}

	err = action.Execute()

	if !assert.NotNil(t, err) {
		return
	}

	assert.IsType(t, tailsafe.ErrAction{}, err)

	t.Run("Type need to be a string", func(t *testing.T) {
		if !assert.NotNil(t, action, "AdapterAction should not be nil") {
			return
		}
		cfg := action.GetConfig()
		err := yaml.Unmarshal([]byte(`{"type": 1}`), cfg)

		err = action.Configure()
		if !assert.NoError(t, err) {
			return
		}

		err = action.Execute()
		if !assert.Error(t, err) {
			return
		}

		assert.IsType(t, tailsafe.ErrAction{}, err)
		assert.ErrorContains(t, err, "type must be in string, not int (1)")
	})
}

func TestAdapterAction_ExecuteWithBadTypeConfig(t *testing.T) {
	step := setuphelper.NewStepTesting()
	action := New(step)

	if !assert.NotNil(t, action, "AdapterAction should not be nil") {
		return
	}

	cfg := action.GetConfig()
	err := yaml.Unmarshal([]byte("bad config"), cfg)

	err = action.Configure()
	if !assert.NoError(t, err) {
		return
	}

	err = action.Execute()
	if !assert.Error(t, err) {
		return
	}

	assert.IsType(t, tailsafe.ErrAction{}, err)
	assert.ErrorContains(t, err, "config is not an object")
}
