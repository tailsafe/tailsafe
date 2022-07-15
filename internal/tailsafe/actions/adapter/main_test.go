package adapteraction

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"github.com/tidwall/gjson"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"
)

type Step struct {
}

func (s Step) GetLogLevel() int {
	return 0
}

func (s *Step) Resolve(path string, data any) any {
	b, err := json.Marshal(data)
	if err != nil {
		return path
	}
	var re = regexp.MustCompile(`(?m){{(.*)}}`)
	result := re.FindAllStringSubmatch(path, -1)
	if len(result) == 0 {
		return path
	}
	value := gjson.Get(string(b), strings.TrimSpace(result[0][1]))
	if value.Type == gjson.String {
		return strings.ReplaceAll(path, result[0][0], fmt.Sprintf("%v", value.String()))
	}
	return value.Value()
}

func (s Step) SetContext(_ context.Context) {
}

func (s Step) SetRuntime(_ tailsafe.RuntimeInterface) {
}

func (s Step) Call() (err error) {
	return
}

func (s Step) GetTitle() string {
	return ""
}

func (s Step) GetUse() string {
	return ""
}

func (s Step) GetContext() context.Context {
	return context.Background()
}

func (s Step) Next() (err tailsafe.ErrActionInterface) {
	return
}

func (s Step) SetCurrent(_ any) {
}

func (s Step) Plugin() tailsafe.ActionInterface {
	return nil
}

func (s Step) Run() error {
	return nil
}

func TestNew(t *testing.T) {
	t.Run("should return a new AdapterAction", func(t *testing.T) {
		instance := New(new(Step))
		if instance == nil {
			t.Error("instance is nil")
		}
	})

	t.Run("should work with object adapter", func(t *testing.T) {
		instance := New(new(Step))

		var data []byte
		data, err := os.ReadFile("./testdata/simple_begin_by_object.yml")
		if err != nil {
			t.Error(err)
			return
		}

		cfg := make(map[string]interface{})
		err = yaml.Unmarshal(data, &cfg)

		if err != nil {
			t.Error(err)
			return
		}

		instance.SetConfig(cfg)
		instance.SetGlobal(map[string]interface{}{
			"test": map[string]interface{}{
				"number":   0,
				"string":   "test",
				"datetime": "2020-01-01T00:00:00Z",
				"boolean":  "true",
			},
		})

		err = instance.Configure()
		if err != nil {
			t.Error(err)
			return
		}

		err = instance.Execute()
		if err != nil {
			t.Error(err)
			return
		}
		enc, err := json.Marshal(instance.GetData())
		if err != nil {
			t.Error(err)
			return
		}
		log.Print(string(enc))
		if string(enc) != `{"boolean":true,"datetime":"2020-01-01T00:00:00Z","number":0,"string":"test","sub_object":{"boolean":true,"datetime":"2020-01-01T00:00:00Z","number":0,"string":"test"}}` {
			t.Error("data is not correct")
		}
	})
	t.Run("should work with array adapter", func(t *testing.T) {
		instance := New(new(Step))

		var data []byte
		data, err := os.ReadFile("./testdata/simple_begin_by_array.yml")
		if err != nil {
			t.Error(err)
			return
		}

		cfg := make(map[string]interface{})
		err = yaml.Unmarshal(data, &cfg)

		if err != nil {
			t.Error(err)
			return
		}

		instance.SetConfig(cfg)
		instance.SetGlobal(map[string]interface{}{
			"global": map[string]interface{}{
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
			},
		})

		err = instance.Configure()
		if err != nil {
			t.Error(err)
			return
		}

		err = instance.Execute()
		if err != nil {
			t.Error(err)
			return
		}

		enc, err := json.Marshal(instance.GetData())
		if err != nil {
			t.Error(err)
			return
		}
		if err != nil {
			t.Error(err)
			return
		}
		if string(enc) != `[{"boolean":true,"datetime":"2020-01-01T00:00:00Z","global_number":999,"number":0,"string":"test","sub_object":{"array":[{"boolean":true,"datetime":"2020-01-01T00:00:00Z","number":0,"string":"test"},{"boolean":false,"datetime":"2020-01-02T00:00:00Z","number":1,"string":"test-1"}],"boolean":true,"datetime":"2020-01-01T00:00:00Z","number":0,"string":"test"}},{"boolean":false,"datetime":"2020-01-02T00:00:00Z","global_number":999,"number":1,"string":"test-1","sub_object":{"array":[{"boolean":true,"datetime":"2020-01-01T00:00:00Z","number":0,"string":"test"},{"boolean":false,"datetime":"2020-01-02T00:00:00Z","number":1,"string":"test-1"}],"boolean":true,"datetime":"2020-01-01T00:00:00Z","number":1,"string":"test"}}]` {
			t.Error("data is not correct", string(enc))
		}
	})

	t.Run("should return a thrown", func(t *testing.T) {
		t.Run("should return a thrown if not supported type", func(t *testing.T) {
			instance := New(new(Step))
			if instance == nil {
				t.Error("instance is nil")
			}

			cfg := make(map[string]interface{})
			cfg["type"] = "??"

			instance.SetConfig(cfg)

			err := instance.Configure()
			if err != nil {
				t.Error(err)
				return
			}

			err = instance.Execute()
			if err == nil {
				t.Error("error is nil")
				return
			}

			castError, ok := err.(tailsafe.ErrActionInterface)
			if !ok {
				t.Error("error is not correct")
			}

			if castError.GetOriginal().Error() != "type `??` is not supported" {
				t.Error("error is not correct")
			}
		})

		t.Run("should return a thrown if invalid type", func(t *testing.T) {
			instance := New(new(Step))
			if instance == nil {
				t.Error("instance is nil")
			}

			cfg := make(map[string]interface{})
			cfg["type"] = 1

			instance.SetConfig(cfg)

			err := instance.Configure()
			if err != nil {
				t.Error(err)
				return
			}

			err = instance.Execute()
			if err == nil {
				t.Error("should return an error")
				return
			}

			castError, ok := err.(tailsafe.ErrActionInterface)
			if !ok {
				t.Error("error is not correct")
			}

			if castError.GetOriginal().Error() != "type must be in string, not int (1)" {
				t.Error("error is not correct")
			}
		})

		t.Run("should return a thrown if data is not a object", func(t *testing.T) {
			instance := New(new(Step))
			if instance == nil {
				t.Error("instance is nil")
			}

			instance.SetConfig(nil)

			err := instance.Configure()
			if err != nil {
				t.Error(err)
				return
			}

			err = instance.Execute()
			if err == nil {
				t.Error("should return an error")
			}

			castError, ok := err.(tailsafe.ErrActionInterface)
			if !ok {
				t.Error("error is not correct")
			}

			log.Print(err)
			if castError.GetOriginal().Error() != "config is not an object" {
				t.Error("error is not correct")
			}
		})
		t.Run("should return a thrown if resolve data doesn't work", func(t *testing.T) {
			instance := New(new(Step))
			if instance == nil {
				t.Error("instance is nil")
			}

			var data []byte
			data, err := os.ReadFile("./testdata/bad_type_in_properties.yml")
			if err != nil {
				t.Error(err)
				return
			}

			cfg := make(map[string]interface{})
			err = yaml.Unmarshal(data, &cfg)

			if err != nil {
				t.Error(err)
				return
			}

			instance.SetConfig(cfg)

			err = instance.Configure()
			if err != nil {
				t.Error(err)
				return
			}

			err = instance.Execute()
			if err == nil {
				t.Error("should return an error")
			}

			castError, ok := err.(tailsafe.ErrActionInterface)
			if !ok {
				t.Error("error is not correct")
			}

			log.Print(err)
			if castError.GetOriginal().Error() != "could not resolve {{ global.array }}" {
				t.Error("error is not correct")
			}
		})
		t.Run("should return a thrown if resolve data is not a slice", func(t *testing.T) {
			instance := New(new(Step))
			if instance == nil {
				t.Error("instance is nil")
			}

			var data []byte
			data, err := os.ReadFile("./testdata/bad_type_in_properties.yml")
			if err != nil {
				t.Error(err)
				return
			}

			cfg := make(map[string]interface{})
			err = yaml.Unmarshal(data, &cfg)

			if err != nil {
				t.Error(err)
				return
			}

			instance.SetConfig(cfg)

			instance.SetGlobal(map[string]interface{}{
				"global": map[string]interface{}{
					"array": 0,
				},
			})

			err = instance.Configure()
			if err != nil {
				t.Error(err)
				return
			}

			err = instance.Execute()
			if err == nil {
				t.Error("should return an error")
			}

			castError, ok := err.(tailsafe.ErrActionInterface)
			if !ok {
				t.Error("error is not correct")
			}

			if castError.GetOriginal().Error() != "float64 is not a slice" {
				t.Error("error is not correct")
			}
		})
	})
}
