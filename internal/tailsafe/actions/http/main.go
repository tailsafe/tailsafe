package httpaction

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type Config struct {
	Headers struct {
		Accept string `json:"accept"`
	} `json:"headers"`
	Method string `json:"method"`
	Path   string `json:"path"`
	URL    string `json:"url"`
}

type HttpAction struct {
	tailsafe.StepInterface
	Config *Config

	global tailsafe.DataInterface
	data   map[string]interface{}
}

func (r *HttpAction) Configure() (err tailsafe.ErrActionInterface) {
	path := r.Resolve(r.Config.Path, r.global)
	if path == nil {
		return tailsafe.CatchStackTrace(r.GetContext(), errors.New("HttpAction: Path is nil"))
	}
	r.Config.Path = path.(string)
	return
}

func (r *HttpAction) Execute() (err tailsafe.ErrActionInterface) {
	requestURL := fmt.Sprintf("%s%s", r.Config.URL, r.Config.Path)
	req, httpErr := http.NewRequest(r.Config.Method, requestURL, nil)
	if httpErr != nil {
		return tailsafe.CatchStackTrace(r.GetContext(), err)
	}

	req.Header.Add("Accept", r.Config.Headers.Accept)

	client := http.Client{}

	// prevent leak
	defer client.CloseIdleConnections()

	// execute request
	res, httpErr := client.Do(req)
	if httpErr != nil {
		return tailsafe.CatchStackTrace(r.GetContext(), httpErr)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	// read response
	body, httpErr := ioutil.ReadAll(res.Body)
	if httpErr != nil {
		return tailsafe.CatchStackTrace(r.GetContext(), httpErr)
	}

	// set data
	r.data["statusCode"] = res.StatusCode
	r.data["headers"] = res.Header

	switch {
	case strings.HasPrefix(r.Config.Headers.Accept, "application/json"):
		var data any
		err := json.Unmarshal(body, &data)
		if err != nil {
			return tailsafe.CatchStackTrace(r.GetContext(), err)
		}
		r.data["body"] = data
		break
	}
	return
}
func (r *HttpAction) GetResult() interface{} {
	return r.data
}
func (r *HttpAction) GetConfig() interface{} {
	if r.Config == nil {
		return &Config{}
	}
	return r.Config
}
func (r *HttpAction) SetConfig(config interface{}) {
	r.Config = config.(*Config)
}
func (r *HttpAction) SetPayload(data tailsafe.DataInterface) {
	// set global data
	r.global = data

	// initialize map
	r.data = make(map[string]interface{})
}
func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(HttpAction)
	p.StepInterface = step
	return p
}
