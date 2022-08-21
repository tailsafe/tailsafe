package httpaction

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Config struct {
	Headers map[string]string `json:"headers"`
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	URL     string            `json:"url"`
	Body    any               `json:"body"`
}

type HttpAction struct {
	tailsafe.StepInterface
	tailsafe.DataInterface
	Config *Config

	data map[string]interface{}
}

func (r *HttpAction) Configure() (err tailsafe.ErrActionInterface) {
	path := r.Resolve(r.Config.Path, r.GetAll())
	if path == nil {
		return tailsafe.CatchStackTrace(r.GetContext(), errors.New("HttpAction: Path is nil"))
	}
	r.Config.Path = path.(string)
	return
}

func (r *HttpAction) Execute() (err tailsafe.ErrActionInterface) {
	requestURL := fmt.Sprintf("%s%s", r.Config.URL, r.Config.Path)

	var payload io.Reader
	switch r.Config.Headers["Content-Type"] {
	case "application/x-www-form-urlencoded":

		values, ok := r.Config.Body.(map[string]any)
		if !ok {
			return
		}

		data := url.Values{}
		for k, v := range values {
			data.Add(k, fmt.Sprintf("%v", r.Resolve(fmt.Sprintf("%v", v), r.GetAll())))
		}

		payload = strings.NewReader(data.Encode())
	}

	req, httpErr := http.NewRequest(r.Config.Method, requestURL, payload)
	if httpErr != nil {
		return tailsafe.CatchStackTrace(r.GetContext(), err)
	}

	for k, v := range r.Config.Headers {
		req.Header.Set(k, r.Resolve(v, r.GetAll()).(string))
	}

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
	case strings.HasPrefix(res.Header.Get("Content-Type"), "application/json"):
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
	return r.Config
}

func (r *HttpAction) SetConfig(config interface{}) {
	r.Config = config.(*Config)
}

func (r *HttpAction) SetPayload(data tailsafe.DataInterface) {
	r.DataInterface = data
}

func New(step tailsafe.StepInterface) tailsafe.ActionInterface {
	p := new(HttpAction)
	p.StepInterface = step
	p.Config = &Config{}
	p.data = make(map[string]interface{})

	return p
}
