package termsAction

import (
	"errors"
	"fmt"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"strconv"
)

type Config struct {
	ActionSetter tailsafe.ActionSetter `json:"action-setter"`
	NoError      bool                  `json:"no-error"`
	Terms        []struct {
		A            string                `json:"a"`
		B            string                `json:"b"`
		Operator     string                `json:"operator"`
		ActionGetter tailsafe.ActionGetter `json:"action-getter"`
	}
}

type If struct {
	tailsafe.StepInterface
	tailsafe.DataInterface

	Config Config
	Result bool
}

func (i *If) Configure() (err tailsafe.ErrActionInterface) {
	if i.Config.Terms == nil {
		return tailsafe.CatchStackTrace(i.GetContext(), errors.New("if: Config has not been set"))
	}
	if len(i.Config.Terms) == 0 {
		return tailsafe.CatchStackTrace(i.GetContext(), errors.New("if: Config is empty"))
	}
	return
}

// Execute executes the action
func (i *If) Execute() (err tailsafe.ErrActionInterface) {
	defer func() {
		i.Set(i.Config.ActionSetter.Key, i.Result, i.Config.ActionSetter.Override)
	}()
	for _, c := range i.Config.Terms {
		switch c.Operator {
		case "eq", "neq":
			switch c.Operator {
			case "eq":
				a := i.Resolve(c.A, i.GetAll())
				b := i.Resolve(c.B, i.GetAll())
				if a == b {
					i.Result = true
					return
				}
				if !i.Config.NoError {
					err = tailsafe.CatchStackTrace(i.GetContext(), errors.New(fmt.Sprintf("Terms: %v is not equal to %v", a, b)))
				}
				break
			case "neq":
				if i.Resolve(c.A, i.GetAll()) != i.Resolve(c.B, i.GetAll()) {
					i.Result = true
				}
				break
			}
			break
		case "gt", "lt", "gte", "lte":
			aV, err := i.toFloat64(i.Resolve(c.A, i.GetAll()))
			if err != nil {
				return tailsafe.CatchStackTrace(i.GetContext(), err)
			}
			bV, err := i.toFloat64(i.Resolve(c.B, i.GetAll()))
			if err != nil {
				return tailsafe.CatchStackTrace(i.GetContext(), err)
			}

			switch c.Operator {
			case "gt":
				if aV > bV {
					i.Result = true
				}
				break
			case "lt":
				if aV < bV {
					i.Result = true
				}
				break
			case "gte":
				if aV >= bV {
					i.Result = true
				}
				break
			case "lte":
				if aV <= bV {
					i.Result = true
				}
				break
			}
		}
	}
	return
}

func (i *If) toFloat64(value any) (float64, error) {
	return strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
}

func (i *If) GetResult() interface{} {
	return i.Result
}

func (i *If) GetConfig() interface{} {
	return &i.Config
}

func (i *If) SetPayload(data tailsafe.DataInterface) {
	i.DataInterface = data
}

func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(If)
	p.StepInterface = step
	p.Config = Config{}
	return p
}
