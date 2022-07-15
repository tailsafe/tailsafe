package If

import (
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"strings"
)

type Config struct {
	Rules []Rule `json:"rules"`
}
type Rule struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	Key   string `json:"key"`
}
type If struct {
	tailsafe.StepInterface
	Config  *Config
	_Global map[string]interface{}
	_Data   any
}

func (r *If) Configure() (err tailsafe.ErrActionInterface) {
	return
}

func (r *If) GetData() any {
	return nil
}

// Execute
// @toDo: add support for multiple rules
// @toDo: extract treatement of rules to a separate function
func (r *If) Execute() tailsafe.ErrActionInterface {
	isContinue := true
	for _, rule := range r.Config.Rules {
		switch rule.Type {
		case "contains":
			key := r.Resolve(rule.Key, r._Data)
			value := r.Resolve(rule.Value, r._Data)
			if key == nil || value == nil {
				break
			}

			// specific case for object
			valueMap, ok := value.(map[string]interface{})
			if ok {
				for _, v := range valueMap {
					keyString, ok := key.(string)
					if !ok {
						break
					}
					if strings.Contains(v.(string), keyString) {
						isContinue = false
						break
					}
				}
				break
			}

			// default case
			keyString, ok := key.(string)
			if !ok {
				break
			}
			valueString, ok := value.(string)
			if !ok {
				break
			}
			if !strings.Contains(valueString, keyString) {
				break
			}
			isContinue = false
		}
	}
	if isContinue {
		return tailsafe.CatchStackTrace(r.GetContext(), &tailsafe.ErrContinue{Message: "continue"})
	}
	// if it matches, we keep the data
	r._Data = r._Global["current"]
	return nil
}
func (r *If) Data() interface{} {
	return r._Data
}
func (r *If) GetConfig() interface{} {
	return &Config{}
}
func (r *If) SetConfig(config interface{}) {
	if config == nil {
		return
	}
	r.Config = config.(*Config)
}
func (r *If) SetGlobal(data map[string]interface{}) {
	r._Global = data
}
func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(If)
	p.StepInterface = step
	return p
}
