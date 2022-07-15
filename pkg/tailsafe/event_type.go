package tailsafe

type ReceiveSignalEvent struct {
	_GenericEvent
	_SignalValue
}

type FileParsedEvent struct {
	_GenericEvent
	_GetTemplate
}

type InitEvent struct {
	_GenericEvent
	_GetVersion
}

type ExitEvent struct {
	_GenericEvent
}
type ExitErrorEvent struct {
	_GenericEvent
	_GetError
}

type ActionBeforeConfigureEvent struct {
	_GenericEvent
	_ActionStepEvent
}

type ActionHasWaitEvent ActionBeforeConfigureEvent

type ActionAfterConfigureEvent struct {
	_GenericEvent
	_ActionStepActionEvent
}

type ActionGenericEvent struct {
	_GenericEvent
	_ActionStepEvent
}

type ActionHasStoringDataEvent struct {
	_GenericEvent
	_ActionStepEvent
	_AnyValue
}

type ActionIsAsyncEvent struct {
	_GenericEvent
	_ActionStepEvent
}

type ActionExitWithErrorEvent struct {
	_GenericEvent
	_ActionStepErrorEvent
	_GetError
	_StageMonitoring
	_IntValue
}

type ActionExitEvent struct {
	_GenericEvent
	_ActionStepEvent
	_StageMonitoring
	_IntValue
}

type ActionStdoutEvent struct {
	_GenericEvent
	_GetStep
	_StringValue
	_IntValue
}

func NewActionHasStoringKeyEvent(step StepInterface) ActionGenericEventInterface {
	instance := &ActionGenericEvent{}
	instance.SetKey(EVENT_ACTION_STORING_KEY)
	instance.SetStep(step)
	return instance
}

func NewActionBeforeConfigureStepEvent(step StepInterface) ActionGenericEventInterface {
	instance := &ActionGenericEvent{}
	instance.SetKey(EVENT_ACTION_BEFORE_CONFIG)
	instance.SetStep(step)
	return instance
}

func NewActionAfterConfigureStepEvent(step StepInterface, action ActionInterface) ActionAfterConfigureEventInterface {
	instance := &ActionAfterConfigureEvent{}
	instance.SetKey(EVENT_ACTION_AFTER_CONFIG)
	instance.SetStep(step)
	instance.SetAction(action)
	return instance
}

func NewFileParsedEvent(template TemplateInterface) FileParsedEventInterface {
	instance := &FileParsedEvent{}
	instance.SetKey(EVENT_FILE_PARSED)
	instance.SetTemplate(template)
	return instance
}

func NewInitEvent(version string, commit string) InitEventInterface {
	instance := &InitEvent{}
	instance.SetKey(EVENT_INIT)
	instance.SetVersion(version)
	instance.SetCommit(commit)
	return instance
}

func NewExitEvent() ExitEventInterface {
	instance := &ExitEvent{}
	instance.SetKey(EVENT_EXIT)
	return instance
}

func NewReceiveSignalEvent(signal string) ReceiveSignalEventInterface {
	instance := &ReceiveSignalEvent{}
	instance.SetKey(EVENT_RECEIVE_SIGNAL)
	instance.SetSignalValue(signal)
	return instance
}

func NewActionIsAsyncEvent(step StepInterface) ActionGenericEventInterface {
	instance := &ActionGenericEvent{}
	instance.SetKey(EVENT_ACTION_IS_ASYNC)
	instance.SetStep(step)
	return instance
}

func NewActionExitWithErrorEvent(step StepInterface, err error, stageMonitoring StageMonitoringInterface, childLevel int) ActionExitWithErrorInterface {
	instance := &ActionExitWithErrorEvent{}
	instance.SetKey(EVENT_ACTION_EXIT_WITH_ERROR)
	instance.SetStep(step)
	instance.SetError(err)
	instance.SetStageMonitoring(stageMonitoring)
	instance.SetIntValue(childLevel)

	return instance
}

func NewExitWithErrorEvent(err error) ActionExitWithErrorInterface {
	instance := &ActionExitWithErrorEvent{}
	instance.SetKey(EVENT_EXIT_WITH_ERROR)
	instance.SetError(err)
	return instance
}

func NewActionExitEvent(step StepInterface, stageMonitoring StageMonitoringInterface, childLevel int) ActionExitEventInterface {
	instance := &ActionExitEvent{}
	instance.SetKey(EVENT_ACTION_EXIT)
	instance.SetStep(step)
	instance.SetStageMonitoring(stageMonitoring)
	instance.SetIntValue(childLevel)
	return instance
}

func NewActionHasStoringDataEvent(step StepInterface, data any) ActionHasStoringDataEventInterface {
	instance := &ActionHasStoringDataEvent{}
	instance.SetKey(EVENT_ACTION_STORING_DATA)
	instance.SetStep(step)
	instance.SetAnyValue(data)
	return instance
}

func NewActionHasWaitEvent(step StepInterface) ActionGenericEventInterface {
	instance := &ActionGenericEvent{}
	instance.SetKey(EVENT_ACTION_HAS_WAIT)
	instance.SetStep(step)
	return instance
}

func NewActionStdoutEvent(step StepInterface, str string, childLevel int) ActionStdoutEventInterface {
	instance := &ActionStdoutEvent{}
	instance.SetKey(EVENT_ACTION_STDOUT)
	instance.SetStep(step)
	instance.SetStringValue(str)
	instance.SetIntValue(childLevel)
	return instance
}

type _GenericEvent struct {
	_Key string
}

func (a *_GenericEvent) Key() string {
	return a._Key
}

func (a *_GenericEvent) SetKey(key string) {
	a._Key = key
}

type _ActionStepEvent struct {
	_GetStep
}

type _ActionStepErrorEvent struct {
	_GetStep
	_GetError
}

type _ActionStepActionEvent struct {
	_GetStep
	_GetAction
}

type _ErrorEvent struct {
	_GenericEvent
	_GetError
}

type _GetError struct {
	Err error
}

func (g *_GetError) GetError() error {
	return g.Err
}
func (g *_GetError) SetError(err error) {
	g.Err = err
}

type _GetTemplate struct {
	Template TemplateInterface
}

func (g _GetTemplate) GetTemplate() TemplateInterface {
	return g.Template
}
func (g *_GetTemplate) SetTemplate(template TemplateInterface) {
	g.Template = template
}

type _GetStep struct {
	Step StepInterface
}

func (g *_GetStep) GetStep() StepInterface {
	return g.Step
}

func (g *_GetStep) SetStep(step StepInterface) {
	g.Step = step
}

type _GetAction struct {
	Action ActionInterface
}

func (g *_GetAction) GetAction() ActionInterface {
	return g.Action
}

func (g *_GetAction) SetAction(action ActionInterface) {
	g.Action = action
}

type _StringValue struct {
	StringValue string
}

func (s *_StringValue) GetStringValue() string {
	return s.StringValue
}

func (s *_StringValue) SetStringValue(stringValue string) {
	s.StringValue = stringValue
}

type _SignalValue struct {
	SignalValue string
}

func (s *_SignalValue) GetSignalValue() string {
	return s.SignalValue
}

func (s *_SignalValue) SetSignalValue(SignalValue string) {
	s.SignalValue = SignalValue
}

type _IntValue struct {
	IntValue int
}

func (i *_IntValue) GetIntValue() int {
	return i.IntValue
}

func (i *_IntValue) SetIntValue(intValue int) {
	i.IntValue = intValue
}

type _StageMonitoring struct {
	StageMonitoring StageMonitoringInterface
}

func (s *_StageMonitoring) GetStageMonitoring() StageMonitoringInterface {
	return s.StageMonitoring
}

func (s *_StageMonitoring) SetStageMonitoring(stageMonitoring StageMonitoringInterface) {
	s.StageMonitoring = stageMonitoring
}

type _AnyValue struct {
	_AnyValue any
}

func (a *_AnyValue) GetAnyValue() any {
	return a._AnyValue
}

func (a *_AnyValue) SetAnyValue(anyValue any) {
	a._AnyValue = anyValue
}

type _GetVersion struct {
	version string
	commit  string
}

func (g *_GetVersion) GetVersion() string {
	return g.version
}
func (g *_GetVersion) GetCommit() string {
	return g.commit
}
func (g *_GetVersion) SetVersion(version string) {
	g.version = version
}
func (g *_GetVersion) SetCommit(commit string) {
	g.commit = commit
}
