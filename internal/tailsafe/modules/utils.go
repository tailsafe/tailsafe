package modules

import (
	"encoding/json"
)

type Utils struct {
}

var UtilsInstance *Utils

func init() {
	UtilsInstance = &Utils{}
}

func GetUtilsModule() *Utils {
	return UtilsInstance
}

// Pretty prints a JSON in a readable format
func (e *Utils) Pretty(v any, level int) any {
	var prefix = "    "
	for i := 0; i < level; i++ {
		prefix += "  "
	}
	b, err := json.MarshalIndent(v, prefix, "  ")
	if err != nil {
		return v
	}
	return string(b)
}

func (e *Utils) Indent(msg string, indentLevel int) string {
	for i := 0; i < indentLevel; i++ {
		msg = "  " + msg
	}
	return msg
}
