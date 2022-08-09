package termsAction

import (
	"errors"
	"fmt"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"strconv"
)

type Config struct {
	A        string `json:"a"`
	B        string `json:"b"`
	Operator string `json:"operator"`
}

type Number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}
type N struct {
	int
	int8
	int16
	int32
	int64
	float32
	float64
}

type If struct {
	tailsafe.StepInterface
	tailsafe.DataInterface

	Config []Config
	Result any
}

func (i *If) Configure() (err tailsafe.ErrActionInterface) {
	if i.Config == nil {
		return tailsafe.CatchStackTrace(i.GetContext(), errors.New("if: Config has not been set"))
	}
	if len(i.Config) == 0 {
		return tailsafe.CatchStackTrace(i.GetContext(), errors.New("if: Config is empty"))
	}
	return
}

// Execute executes the action
func (i *If) Execute() tailsafe.ErrActionInterface {
	defer func() {
		i.Set(tailsafe.RETURN, i.Result)
	}()
	for _, c := range i.Config {
		switch c.Operator {
		case "==", "!=":
			switch c.Operator {
			case "==":
				if i.Resolve(c.A, map[string]any{"this": i.Get(tailsafe.THIS)}) == i.Resolve(c.B, map[string]any{"this": i.Get(tailsafe.THIS)}) {
					i.Result = true
				}
				break
			case "!=":
				if i.Resolve(c.A, map[string]any{"this": i.Get(tailsafe.THIS)}) != i.Resolve(c.B, map[string]any{"this": i.Get(tailsafe.THIS)}) {
					i.Result = true
				}
				break
			}
			break
		case ">", "<", ">=", "<=":
			aV, err := i.toFloat64(i.Resolve(c.A, map[string]any{"this": i.Get(tailsafe.THIS)}))
			if err != nil {
				return tailsafe.CatchStackTrace(i.GetContext(), err)
			}
			bV, err := i.toFloat64(i.Resolve(c.B, map[string]any{"this": i.Get(tailsafe.THIS)}))
			if err != nil {
				return tailsafe.CatchStackTrace(i.GetContext(), err)
			}

			switch c.Operator {
			case ">":
				if aV > bV {
					i.Result = true
				}
				break
			case "<":
				if aV < bV {
					i.Result = true
				}
				break
			case ">=":
				if aV >= bV {
					i.Result = true
				}
				break
			case "<=":
				if aV <= bV {
					i.Result = true
				}
				break
			}
		}
	}
	return nil
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
	p.Config = []Config{}
	return p
}
