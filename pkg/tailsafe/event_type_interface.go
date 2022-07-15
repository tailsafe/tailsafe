package tailsafe

const (
	/* GLOBAL */

	EVENT_INIT            = "init"
	EVENT_RECEIVE_SIGNAL  = "receive_signal"
	EVENT_EXIT            = "exit"
	EVENT_EXIT_WITH_ERROR = "exit_with_error"
	EVENT_FILE_PARSED     = "file_parsed"

	/* ACTION EVENT */

	EVENT_ACTION_HAS_WAIT        = "action_has_wait"
	EVENT_ACTION_IS_ASYNC        = "action_is_async"
	EVENT_ACTION_STORING_KEY     = "action_storing_key"
	EVENT_ACTION_BEFORE_CONFIG   = "action_before_config"
	EVENT_ACTION_AFTER_CONFIG    = "action_after_config"
	EVENT_ACTION_EXIT_WITH_ERROR = "action_exit_with_error"
	EVENT_ACTION_STORING_DATA    = "action_storing_data"
	EVENT_ACTION_EXIT            = "action_exit"
	EVENT_ACTION_STDOUT          = "action_stout"
)

type EventInterface interface {
	Key() string
}

type ReceiveSignalEventInterface interface {
	EventInterface
	GetSignalValue
}

type FileParsedEventInterface interface {
	EventInterface
	GetTemplate
}

type InitEventInterface interface {
	EventInterface
	GetVersion
}

type ExitEventInterface interface {
	EventInterface
}

type ActionGenericEventInterface interface {
	EventInterface
	GetStep
}

type ActionHasStoringKeyEventInterface interface {
	EventInterface
	GetStep
}

type ActionHasWaitEventInterface ActionHasStoringKeyEventInterface

type ActionHasStoringDataEventInterface interface {
	EventInterface
	GetStep
	GetAnyValue
}

type ActionAfterConfigureEventInterface interface {
	EventInterface
	GetStep
	GetActionEvent
}

type ActionBeforeConfigureEventInterface interface {
	EventInterface
	GetStep
}

type ActionIsAsyncInterface interface {
	EventInterface
	GetStep
}

type ActionExitWithErrorInterface interface {
	EventInterface
	GetStep
	GetError
	GetStageMonitoring
	GetIntValue
}

type ErrorEventInterface interface {
	EventInterface
	GetError
}

type ActionExitEventInterface interface {
	EventInterface
	GetStep
	GetStageMonitoring
	GetIntValue
}

type ActionStdoutEventInterface interface {
	EventInterface
	GetStep
	GetStringValue
	GetIntValue
}

type GetStringValue interface {
	GetStringValue() string
	SetStringValue(string)
}

type GetSignalValue interface {
	GetSignalValue() string
	SetSignalValue(string)
}

type GetAnyValue interface {
	GetAnyValue() any
	SetAnyValue(any)
}

type GetIntValue interface {
	GetIntValue() int
	SetIntValue(int)
}

type GetStep interface {
	GetStep() StepInterface
	SetStep(StepInterface)
}
type GetError interface {
	GetError() error
	SetError(error)
}
type GetActionEvent interface {
	GetAction() ActionInterface
	SetAction(ActionInterface)
}

type GetTemplate interface {
	GetTemplate() TemplateInterface
	SetTemplate(TemplateInterface)
}

type GetStageMonitoring interface {
	GetStageMonitoring() StageMonitoringInterface
	SetStageMonitoring(StageMonitoringInterface)
}

type GetVersion interface {
	GetVersion() string
	SetVersion(string)

	GetCommit() string
	SetCommit(string)
}
