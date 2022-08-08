package modules

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUtilsModule(t *testing.T) {
	assert.NotNil(t, GetUtilsModule())
}

func TestUtils_GetAppActionDir(t *testing.T) {
	assert.Contains(t, GetUtilsModule().GetAppActionDir(), "action")
}

func TestUtils_Indent(t *testing.T) {
	assert.Equal(t, "  ", GetUtilsModule().Indent("", 1))
}

func TestUtils_GetAppTemplateDir(t *testing.T) {
	assert.Contains(t, GetUtilsModule().GetAppTemplateDir(), "template")
}

func TestUtils_GetAppDir(t *testing.T) {
	assert.Contains(t, GetUtilsModule().GetAppDir(), "/.tailsafe")
}

func TestUtils_Pretty(t *testing.T) {
	assert.Equal(t, "{\n        \"test\": \"test\"\n      }", GetUtilsModule().Pretty(map[string]string{"test": "test"}, 1))
}
func TestUtils_PrettyWithFail(t *testing.T) {
	v := make(chan int)
	assert.Equal(t, v, GetUtilsModule().Pretty(v, 1))
}
