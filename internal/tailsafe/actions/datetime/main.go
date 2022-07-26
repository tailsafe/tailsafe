package datetime

import (
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"time"
)

type Config struct {
	Format string `json:"format"`
	Add    struct {
		Year  int `json:"year"`
		Month int `json:"month"`
		Day   int `json:"day"`
	} `json:"add"`
}

type DateTime struct {
	tailsafe.StepInterface
	Config  *Config
	_Global tailsafe.DataInterface
	_Data   string
}

func (r *DateTime) Configure() (err tailsafe.ErrActionInterface) {
	return
}

func (r *DateTime) GetResult() any {
	return nil
}

func (r *DateTime) Execute() (err tailsafe.ErrActionInterface) {
	date := time.Now()
	date = date.AddDate(r.Config.Add.Year, r.Config.Add.Month, r.Config.Add.Day)

	r._Data = date.String()
	if r.Config.Format != "" {
		r._Data = date.Format(r.Config.Format)
	}
	return
}
func (r *DateTime) GetConfig() interface{} {
	return &Config{}
}
func (r *DateTime) SetConfig(config interface{}) {
	r.Config = config.(*Config)
}
func (r *DateTime) SetPayload(data tailsafe.DataInterface) {
	r._Global = data
}
func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(DateTime)
	p.StepInterface = step
	return p
}
