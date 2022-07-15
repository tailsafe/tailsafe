package modules

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestModule struct {
}

func (t *TestModule) GetName() string {
	return "test"
}

type TestModuleInterface interface {
	GetName() string
}

func TestRegister(t *testing.T) {
	Register("test", &TestModule{})
}

func TestRegisterWithError(t *testing.T) {
	assert.Panics(t, func() { Register("test", &TestModule{}) })
}

func TestGet(t *testing.T) {
	module := Get[TestModuleInterface]("test")
	assert.Equalf(t, "test", module.GetName(), "Expected module name to be 'test'")
}

func TestGetWithError(t *testing.T) {
	assert.Panics(t, func() { Get[TestModuleInterface]("test_2") })
}

func TestRequires(t *testing.T) {
	err := Requires([]string{"test"})
	assert.Nil(t, err)
}

func TestRequiresWithError(t *testing.T) {
	err := Requires([]string{"test_2"})
	assert.NotNilf(t, err, "Expected error when getting module 'test_2'")
}

func TestReset(t *testing.T) {
	Reset()
	assert.Equalf(t, 0, len(modules), "Expected modules to be empty")
}

func BenchmarkGet(b *testing.B) {
	Register("test", &TestModule{})

	for i := 0; i < b.N; i++ {
		_ = Get[TestModuleInterface]("test")
	}
}
