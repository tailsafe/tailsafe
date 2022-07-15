package replace

import (
	"fmt"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"regexp"
)

type Config struct {
	Rules []struct {
		Old string `json:"old"`
		New string `json:"new"`
		Key string `json:"key"`
	} `json:"rules"`
}
type Replace struct {
	tailsafe.StepInterface
	Config  *Config
	_Global map[string]interface{}
	_Data   []map[string]interface{}
}

func (r *Replace) Configure() (err tailsafe.ErrActionInterface) {
	return
}

func (r *Replace) GetData() any {
	return nil
}

func (r *Replace) Execute() (err tailsafe.ErrActionInterface) {
	_, ok := r._Global["current"].(map[string]interface{})
	if !ok {
		return
	}
	for _, rule := range r.Config.Rules {
		var re = regexp.MustCompile(`(?m)^.+`)
		var str = fmt.Sprintf("%v", r._Global["current"].(map[string]interface{})[rule.Key])
		r._Global["current"].(map[string]interface{})[rule.Key] = re.ReplaceAllString(str, rule.New)
	}
	return
}
func (r *Replace) Data() interface{} {
	return r._Global["current"]
}
func (r *Replace) GetConfig() interface{} {
	return &Config{}
}
func (r *Replace) SetConfig(config interface{}) {
	r.Config = config.(*Config)
}
func (r *Replace) SetGlobal(data map[string]interface{}) {
	r._Global = data
}
func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(Replace)
	p.StepInterface = step
	return p
}
