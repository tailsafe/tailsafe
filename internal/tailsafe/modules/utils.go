package modules

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Utils struct {
	appDir         string
	appDirAction   string
	appDirTemplate string
}

var UtilsInstance *Utils

func init() {
	UtilsInstance = &Utils{}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	UtilsInstance.appDir = filepath.Join(home, ".tailsafe")
	UtilsInstance.appDirAction = filepath.Join(home, ".tailsafe", "actions")
	UtilsInstance.appDirTemplate = filepath.Join(home, ".tailsafe", "template")

	err = os.MkdirAll(UtilsInstance.appDir, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	err = os.MkdirAll(UtilsInstance.appDirAction, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	err = os.MkdirAll(UtilsInstance.appDirTemplate, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
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

func (e *Utils) GetAppDir() string {
	return e.appDir
}

func (e *Utils) GetAppActionDir() string {
	return e.appDir
}

func (e *Utils) GetAppTemplateDir() string {
	return e.appDirTemplate
}
