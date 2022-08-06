package main

import (
	"github.com/tailsafe/tailsafe/internal/tailsafe/modules"
	tailsafeInterface "github.com/tailsafe/tailsafe/pkg/tailsafe"
	"time"
)

func init() {
	// Init timer
	var timeStart = time.Now()
	var last = ""

	// Listen engine events.
	logger := modules.GetLoggerModule()
	modules.GetEventsModule().Subscribe(func(event tailsafeInterface.EventInterface) {
		switch e := event.(type) {
		case tailsafeInterface.InitEventInterface:
			payload := logger.
				NewPayload().
				SetNamespace(tailsafeInterface.NAMESPACE_DEFAULT).
				SetLevel(tailsafeInterface.LOG_INFO).
				SetMessage("\n")

			modules.
				GetLoggerModule().
				Log(payload)

			modules.
				GetLoggerModule().
				Log(payload.SetMessage("🚀 Launching tailsafe-cli - v%s #%s\n\n").SetArgs("1.0.0", "6ca5590"))
			break
		case tailsafeInterface.ReceiveSignalEventInterface:
			payload := logger.
				NewPayload().
				SetNamespace(tailsafeInterface.NAMESPACE_DEFAULT).
				SetLevel(tailsafeInterface.LOG_INFO).
				SetMessage("\n \u001B[33m>>> Received %s signal, closing all actions...\u001B[0m\n\n").
				SetArgs(e.GetSignalValue())

			modules.
				GetLoggerModule().
				Log(payload)
			break

		case tailsafeInterface.FileParsedEventInterface:
			payload := logger.
				NewPayload().
				SetNamespace(tailsafeInterface.NAMESPACE_WORKFLOW).
				SetLevel(tailsafeInterface.LOG_INFO)

			modules.GetLoggerModule().Log(payload.SetMessage("\n"))
			modules.GetLoggerModule().Log(payload.SetMessage("\u001B[36mTitle       :\u001B[0m %s").SetArgs(e.GetTemplate().GetTitle()))
			modules.GetLoggerModule().Log(payload.SetMessage("\u001B[36mDescription :\u001B[0m %s").SetArgs(e.GetTemplate().GetDescription()))
			modules.GetLoggerModule().Log(payload.SetMessage("\u001B[36mMaintainer  :\u001B[0m %s").SetArgs(e.GetTemplate().GetMaintainer()))
			modules.GetLoggerModule().Log(payload.SetMessage("\u001B[36mRevision    :\u001B[0m %d").SetArgs(e.GetTemplate().GetRevision()))
			modules.GetLoggerModule().Log(payload.SetMessage("\n"))
			modules.GetLoggerModule().Log(payload.SetMessage("\u001B[32m>> Launching of the stages ...\033[0m"))
			break

		case tailsafeInterface.ActionAfterConfigureEventInterface:
			payload := logger.
				NewPayload().
				SetNamespace(tailsafeInterface.NAMESPACE_WORKFLOW).
				SetLevel(e.GetStep().GetLogLevel())

			modules.GetLoggerModule().Log(payload.SetMessage("\u001B[36m  ↳ Config:\u001B[0m %v").SetArgs(modules.Get[tailsafeInterface.Utils]("Utils").Pretty(e.GetAction().GetConfig(), e.GetStep().GetLogLevel())))
			break
		case tailsafeInterface.ActionExitWithErrorInterface:
			payload := logger.
				NewPayload()

			if e.GetStep() != nil {
				payload.
					SetNamespace(tailsafeInterface.NAMESPACE_WORKFLOW).
					SetLevel(e.GetStep().GetLogLevel())

				modules.GetLoggerModule().Log(payload.SetMessage("  ↳ Stage %d \u001B[31maborded\u001B[0m after %s").SetArgs(e.GetStageMonitoring().GetStage(), e.GetStageMonitoring().GetStageDuration()))
			}

			payload.
				SetNamespace(tailsafeInterface.NAMESPACE_DEFAULT).
				SetLevel(tailsafeInterface.LOG_INFO).
				SetMessage("\n\u001B[31m>> An error has occurred\u001B[0m")

			modules.
				GetLoggerModule().
				Log(payload)

			modules.
				GetLoggerModule().
				Log(payload.SetMessage("\n\u001B[31mError :\u001B[0m %s").SetArgs(e.GetError()))

			switch t := e.GetError().(type) {
			case tailsafeInterface.ErrActionInterface:
				modules.
					GetLoggerModule().
					Log(payload.SetMessage("\n\u001B[31mStacktrace :\u001B[0m %s").SetLevel(tailsafeInterface.LOG_VERBOSE).SetArgs(t.GetStackTrace()))
			}
			modules.
				GetLoggerModule().
				Log(payload.SetMessage("\n").SetLevel(tailsafeInterface.LOG_INFO))
			break

		case tailsafeInterface.ActionExitEventInterface:
			payload := logger.
				NewPayload().
				SetNamespace(tailsafeInterface.NAMESPACE_WORKFLOW).
				SetLevel(e.GetStep().GetLogLevel())

			modules.GetLoggerModule().Log(payload.SetMessage(modules.GetUtilsModule().Indent("  ↳ Stage %d \u001B[32mcompleted\u001B[0m in %s", e.GetIntValue())).SetArgs(e.GetStageMonitoring().GetStage(), e.GetStageMonitoring().GetStageDuration()))
			break

		case tailsafeInterface.ActionStdoutEventInterface:
			if last != e.GetStep().GetKey() {
				modules.
					GetLoggerModule().
					Log(logger.NewPayload().SetNamespace(tailsafeInterface.NAMESPACE_DEFAULT).SetLevel(e.GetStep().GetLogLevel()).SetMessage(modules.GetUtilsModule().Indent("\n \u001B[33m>>> Stdout from `%s`\u001B[0m", e.GetIntValue())).SetArgs(e.GetStep().GetTitle()))

				last = e.GetStep().GetKey()
			}

			payload := logger.
				NewPayload().
				SetNamespace(tailsafeInterface.NAMESPACE_DEFAULT).
				SetLevel(e.GetStep().GetLogLevel()).
				SetMessage("%s").
				SetArgs(modules.GetUtilsModule().Indent(e.GetStringValue(), e.GetIntValue()))

			modules.
				GetLoggerModule().
				Log(payload)
			break

		case tailsafeInterface.ActionGenericEventInterface:
			switch e.Key() {
			case tailsafeInterface.EVENT_ACTION_STORING_KEY:
				payload := logger.
					NewPayload().
					SetNamespace(tailsafeInterface.NAMESPACE_WORKFLOW).
					SetLevel(e.GetStep().GetLogLevel())

				modules.GetLoggerModule().Log(payload.SetMessage("\u001B[36m  ↳ Storing key :\u001B[0m %s").SetArgs(e.GetStep().GetKey()))

				break
			case tailsafeInterface.EVENT_ACTION_IS_ASYNC:
				payload := logger.
					NewPayload().
					SetNamespace(tailsafeInterface.NAMESPACE_WORKFLOW).
					SetLevel(e.GetStep().GetLogLevel()).
					SetMessage("\u001B[36m  ↳ Async :\u001B[0m %v").
					SetArgs(e.GetStep().IsAsync())

				modules.
					GetLoggerModule().
					Log(payload)

				break
			case tailsafeInterface.EVENT_ACTION_BEFORE_CONFIG:
				payload := logger.
					NewPayload().
					SetNamespace(tailsafeInterface.NAMESPACE_WORKFLOW).
					SetLevel(e.GetStep().GetLogLevel())

				modules.GetLoggerModule().Log(payload.SetMessage("\n"))
				modules.GetLoggerModule().Log(payload.SetMessage("\u001B[36m→ Stage %d :\033[0m %s").SetArgs(e.GetStep().GetEngine().GetStage(), e.GetStep().GetTitle()))
				modules.GetLoggerModule().Log(payload.SetMessage("\u001B[36m  ↳ Action :\u001B[0m %s").SetArgs(e.GetStep().GetUse()))
				break

			case tailsafeInterface.EVENT_ACTION_HAS_WAIT:
				payload := logger.
					NewPayload().
					SetNamespace(tailsafeInterface.NAMESPACE_WORKFLOW).
					SetLevel(e.GetStep().GetLogLevel()).
					SetMessage("\u001B[36m  ↳ Action has wait :\u001B[0m %v").
					SetArgs(e.GetStep().GetWait())

				modules.
					GetLoggerModule().
					Log(payload)
				break
			}
			break
		case tailsafeInterface.ExitEventInterface:
			payload := logger.
				NewPayload().
				SetNamespace(tailsafeInterface.NAMESPACE_DEFAULT).
				SetLevel(tailsafeInterface.LOG_INFO).
				SetMessage("\n")

			modules.
				GetLoggerModule().
				Log(payload)

			modules.
				GetLoggerModule().
				Log(payload.SetMessage("\u001B[32m>> Successfully completed in %s️.\u001B[0m ❤️Thanks for use ❤").SetArgs(time.Since(timeStart)))
		}
	})
}
