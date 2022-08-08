package modules

import (
	"fmt"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"golang.org/x/exp/slices"
	"log"
	"regexp"
)

type Logger struct {
	stageLevel int
	logColor   bool
	namespaces []string
	verbose    bool
}

type LoggerPayload struct {
	Namespace string
	Level     int
	Message   string
	Args      []any
}

func (l *LoggerPayload) GetNamespace() string {
	return l.Namespace
}

func (l *LoggerPayload) GetLevel() int {
	return l.Level
}

func (l *LoggerPayload) GetMessage() string {
	return l.Message
}

func (l *LoggerPayload) GetArgs() []any {
	return l.Args
}

func (l *LoggerPayload) SetNamespace(namespace string) tailsafe.LoggerPayload {
	l.Namespace = namespace
	return l
}

func (l *LoggerPayload) SetLevel(i int) tailsafe.LoggerPayload {
	l.Level = i
	return l
}

func (l *LoggerPayload) SetMessage(message string) tailsafe.LoggerPayload {
	// reset args
	l.SetArgs()

	// set message
	l.Message = message
	return l
}

func (l *LoggerPayload) SetArgs(args ...any) tailsafe.LoggerPayload {
	l.Args = args
	return l
}

var loggerInstance *Logger

func init() {
	// create the logger instance with default values
	loggerInstance = &Logger{}
	loggerInstance.SetStageLevel(tailsafe.LOG_INFO)
	loggerInstance.SetLogColor(true)
	loggerInstance.namespaces = append(loggerInstance.namespaces, tailsafe.NAMESPACE_DEFAULT)
}

func GetLoggerModule() *Logger {
	return loggerInstance
}

func (l *Logger) NewPayload() tailsafe.LoggerPayload {
	return new(LoggerPayload)
}

func (l *Logger) SetStageLevel(level int) *Logger {
	l.stageLevel = level
	return l
}

func (l *Logger) SetLogColor(color bool) *Logger {
	l.logColor = color
	return l
}

func (l *Logger) AddNamespace(namespace string) *Logger {
	l.namespaces = append(l.namespaces, namespace)
	return l
}

func (l *Logger) SetVerbose(verbose bool) *Logger {
	l.verbose = verbose
	return l
}

func (l *Logger) GetStageLevel() int {
	return l.stageLevel
}

func (l *Logger) GetLogColor() bool {
	return l.logColor
}

func (l *Logger) GetNamespaces() []string {
	return l.namespaces
}

func (l *Logger) GetVerbose() bool {
	return l.verbose
}

// Log logs a message with the given level
func (l *Logger) Log(payload tailsafe.LoggerPayload) {
	// if namespace contains the namespace, then log the message
	if !slices.Contains(l.GetNamespaces(), payload.GetNamespace()) {
		return
	}
	// if log disabled
	if payload.GetLevel() == tailsafe.LOG_NONE {
		return
	}
	// if log verbose is disabled
	if payload.GetLevel() > tailsafe.LOG_INFO && !l.GetVerbose() {
		return
	}
	// remove color strip
	if !l.GetLogColor() {
		var re = regexp.MustCompile(`\x1b\[[\d;]*m`)
		payload.SetMessage(re.ReplaceAllString(payload.GetMessage(), ""))
	}
	// print the message
	log.Print(fmt.Sprintf(payload.GetMessage(), payload.GetArgs()...))
}
