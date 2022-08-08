package resolver

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	example = `
{
    "message": "hello",
    "messages": [
        {
            "id": 1
        },
        {
            "id": 2,
            "label": [
                "open",
                "closed"
            ]
        }
    ],
    "state": {
        "label": "open"
    }
}`
	exampleList = `
[
	{"id": 1},
	{"id": 2}
]`
)

func BenchmarkGetSimple(b *testing.B) {
	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(example), &data)
	if err != nil {
		return
	}

	for i := 0; i < b.N; i++ {
		_ = Get("message", data)
	}
}
func BenchmarkGetRecursive(b *testing.B) {
	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(example), &data)
	if err != nil {
		return
	}

	for i := 0; i < b.N; i++ {
		_ = Get("state.label", data)
	}
}

func TestGetSimple(t *testing.T) {
	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(example), &data)
	if !assert.Nil(t, err) {
		return
	}

	message := Get("message", data)
	assert.Equal(t, "hello", message)
	t.Logf("value=%s", message)
}
func TestGetSimpleNotFoundField(t *testing.T) {
	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(example), &data)
	if !assert.Nil(t, err) {
		return
	}

	notfound := Get("nofound", data)
	assert.Equal(t, nil, notfound)
}
func TestGetSimpleIncorrectField(t *testing.T) {
	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(example), &data)
	if !assert.Nil(t, err) {
		return
	}

	notfound := Get("", data)
	assert.Equal(t, nil, notfound)
}
func TestGetRecursive(t *testing.T) {
	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(example), &data)
	if !assert.Nil(t, err) {
		return
	}
	message := Get("state.label", data)
	assert.Equal(t, "open", message)
	t.Logf("value=%s", message)
}
func TestGetRecursiveArray(t *testing.T) {
	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(example), &data)
	if !assert.Nil(t, err) {
		return
	}
	message := Get("messages.1.label.0", data)
	assert.Equal(t, "open", message)
	t.Logf("value=%s", message)
}
func TestGetRecursiveArrayWithNoIndex(t *testing.T) {
	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(example), &data)
	if !assert.Nil(t, err) {
		return
	}
	message := Get("messages.1.label.unknow", data)
	assert.Equal(t, nil, message)
}
func TestGetObjectChildArrayKey(t *testing.T) {
	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(example), &data)
	if !assert.Nil(t, err) {
		return
	}
	id := Get("messages.0.id", data)
	assert.Equal(t, float64(1), id)
	t.Logf("value=%v", id)
}
func TestGetArrayKey(t *testing.T) {
	var data []interface{}
	err := json.Unmarshal([]byte(exampleList), &data)
	if !assert.Nil(t, err) {
		return
	}
	id := Get("0.id", data)
	assert.Equal(t, float64(1), id)
	t.Logf("value=%v", id)
}
func TestGetArrayKeyButNotExist(t *testing.T) {
	var data []interface{}
	err := json.Unmarshal([]byte(exampleList), &data)
	if !assert.Nil(t, err) {
		return
	}
	id := Get("10.id", data)
	assert.Equal(t, nil, id)
}
func TestGetArrayKeyIncorrectField(t *testing.T) {
	var data []interface{}
	err := json.Unmarshal([]byte(exampleList), &data)
	if !assert.Nil(t, err) {
		return
	}
	id := Get("", data)
	assert.Equal(t, nil, id)
}
func TestGetIncorrectData(t *testing.T) {
	id := Get("unknow", nil)
	assert.Equal(t, nil, id)
}
