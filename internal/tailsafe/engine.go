package tailsafe

import (
	"context"
	"errors"
	"fmt"
	"github.com/tailsafe/tailsafe/internal/tailsafe/data"
	"github.com/tailsafe/tailsafe/internal/tailsafe/modules"
	"github.com/tailsafe/tailsafe/internal/tailsafe/versions"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"sync"
	"syscall"
)

const (
	exitCodeErr       = 1
	exitCodeInterrupt = 2
)

type Engine struct {
	tailsafe.DataInterface

	ctx      context.Context
	path     string
	pathData string
	env      string

	// data processing
	data     map[string]any
	mockData map[string]any
	template versions.TemplateInterface
	mu       *sync.Mutex

	// Log process
	childLevel int
	stageLevel int
	logColor   bool

	modules map[string]any
}

func (e *Engine) GetMockDataByKey(key string) any {
	data, ok := e.GetMockData()[key]
	if !ok {
		return nil
	}
	return data
}

// NewStage increments the stage level
func (e *Engine) NewStage() {
	e.stageLevel++
}

// GetStage returns the current stage level
func (e *Engine) GetStage() int {
	return e.stageLevel
}

func (e *Engine) GetChildLevel() int {
	return e.childLevel
}

func (e *Engine) EntrySubStage() {
	e.childLevel++
}

func (e *Engine) ExitSubStage() {
	e.childLevel--
}

// Context returns the context of the payload
func (e *Engine) Context() context.Context {
	return e.ctx
}

// SetContext sets the context of the payload
func (e *Engine) SetContext(ctx context.Context) *Engine {
	e.ctx = ctx
	return e
}

// New creates a new payload
func New() tailsafe.EngineInterface {

	// disable flag for logging
	log.SetFlags(0)

	// new instance of tailsafe-cli
	p := new(Engine)
	// set default context
	p.ctx = context.Background()
	// initialize default map
	p.data = make(map[string]any)
	// initialize default mock data
	p.mockData = make(map[string]any)
	// initialize default mutex
	p.mu = new(sync.Mutex)
	// initialize default payload
	p.DataInterface = data.NewPayload()

	var cancel context.CancelFunc
	_, cancel = context.WithCancel(p.ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case v := <-c:
			modules.Get[tailsafe.EventsInterface]("Events").Trigger(tailsafe.NewReceiveSignalEvent(v.String()))
			cancel()
		}
		// if second try, hard exit
		<-c
		os.Exit(exitCodeInterrupt)
	}()

	return p
}

// SetPath sets the path of the payload
func (e *Engine) SetPath(path string) tailsafe.EngineInterface {
	e.path = path
	return e
}

// SetDataPath sets the path of the payload
func (e *Engine) SetDataPath(path string) tailsafe.EngineInterface {
	e.pathData = path
	return e
}

// SetEnv sets the path of the payload
func (e *Engine) SetEnv(env string) tailsafe.EngineInterface {
	e.env = env
	return e
}

// GetPath returns the path of the payload
func (e *Engine) GetPath() string {
	return e.path
}

// GetPathData returns the path for the mock data
func (e *Engine) GetPathData() string {
	return e.pathData
}

// GetMutex returns the mutex of the payload
func (e *Engine) GetMutex() *sync.Mutex {
	return e.mu
}

// GetData returns the data of the payload
func (e *Engine) GetData() map[string]any {
	return e.data
}

// GetMockData returns the mock data of the payload
func (e *Engine) GetMockData() map[string]any {
	return e.mockData
}

// Parse checks if the payload is valid
func (e *Engine) Parse() (err error) {
	newPath := fmt.Sprintf("%s/%s.yml", modules.GetUtilsModule().GetAppTemplateDir(), e.GetPath())
	if _, err := os.Stat(newPath); err == nil {
		e.SetPath(newPath)
	}

	var data []byte
	data, err = os.ReadFile(e.GetPath())
	if err != nil {
		return
	}
	var template map[string]any
	err = yaml.Unmarshal(data, &template)
	if err != nil {
		return
	}

	// Check if the template contains the required fields
	v, ok := template["version"]
	if !ok {
		return errors.New("version is missing")
	}
	// Check if versions is a string
	vString, ok := v.(string)
	if !ok {
		return errors.New("version is not a string")
	}
	// Check if the versions is supported
	e.template, err = versions.GetTemplate(data, vString)
	if err != nil {
		return
	}

	// Check if user as specified a path to the data
	// @TODO: need refactoring override short path
	if e.GetPathData() != "" {
		newPathData := fmt.Sprintf("%s/%s.yml", modules.GetUtilsModule().GetAppTemplateDir(), e.GetPathData())
		if _, err := os.Stat(newPathData); err == nil {
			e.SetDataPath(newPathData)
		}

		var mockData []byte
		mockData, err = os.ReadFile(e.GetPathData())
		if err != nil {
			return
		}
		var dataTemplate map[string]any
		err = yaml.Unmarshal(mockData, &dataTemplate)
		if err != nil {
			return
		}
		// set the mock data
		for k, v := range dataTemplate {
			e.mockData[k] = v
		}
	}
	return
}

// Run executes the payload
func (e *Engine) Run() {
	var err error
	defer func() {
		if err == nil {
			modules.Get[tailsafe.EventsInterface]("Events").Trigger(tailsafe.NewExitEvent())
			return
		}
		modules.Get[tailsafe.EventsInterface]("Events").Trigger(tailsafe.NewExitWithErrorEvent(err))
		os.Exit(exitCodeErr)
	}()

	err = modules.Requires([]string{
		"Utils",
		"Events",
		"AsyncQueue",
	})

	if err != nil {
		return
	}

	modules.Get[tailsafe.EventsInterface]("Events").Trigger(tailsafe.NewInitEvent("dev", "dev"))

	// Check if the payload is valid
	err = e.Parse()
	if err != nil {
		return
	}

	// file is parsed, now we can start the execution
	modules.Get[tailsafe.EventsInterface]("Events").Trigger(tailsafe.NewFileParsedEvent(e.template))

	// Get dependencies
	requires := e.template.GetDependencies()
	for _, require := range requires {
		// not use autoload because is already compiled with engine.
		if strings.HasPrefix(require, "internal") {
			continue
		}
		// add compile task
		if _, err := os.Stat(require); !os.IsNotExist(err) {
			split := strings.Split(require, "/")
			name := split[len(split)-1]

			step := e.template.NewStep()
			step.SetUse("internal/exec")
			step.SetTitle("[DEV] Build Action")
			step.SetConfig(map[string]any{
				"command": []string{
					"go",
					"build",
					"-buildmode=plugin",
					"-o",
					fmt.Sprintf("%s/%s@dev.so", modules.GetUtilsModule().GetAppActionDir(), name),
					".",
				}, "path": require})

			e.template.InjectPreStep([]tailsafe.StepInterface{step})
			continue
		}
	}

	// Validate arguments
	args, err := e.template.SetEnv(e.env)
	if err != nil {
		return
	}

	// @todo search a better idea for key wording please :D
	e.Set("SYSTEM_ARGS", args, true)

	// convert the data to the correct type for iteration
	steps, err := InterfaceSlice[tailsafe.StepInterface](e.template.GetSteps())
	if err != nil {
		return
	}
	// Execute the actions into steps
	for _, step := range steps {
		// set the context into the action
		step.SetContext(e.ctx)
		step.SetEngine(e)

		// call the action
		err = step.Call()
		if err != nil {
			return
		}

		if step.HasFailed() {
			return
		}
	}

	// Wait all the actions to finish
	modules.GetAsyncQueue().WaitAll()
	return
}

// ExtractGlobal extract required data from the global context
func (e *Engine) ExtractGlobal(required []string) map[string]any {
	global := make(map[string]any)
	for _, v := range required {
		data, ok := e.GetData()[v]
		if !ok {
			continue
		}
		global[v] = data
	}
	return global
}

// InterfaceSlice helps to convert a slice of interfaces to a slice of a specific type
// returns an error if the type is not supported
// use generic type to avoid type checking
func InterfaceSlice[t any](slice any) (data []t, err error) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		err = errors.New(fmt.Sprintf("%s is not a slice", reflect.TypeOf(slice)))
		return
	}
	if s.IsNil() {
		return
	}
	data = make([]t, s.Len())
	for i := 0; i < s.Len(); i++ {
		data[i] = s.Index(i).Interface().(t)
	}
	return
}
