package selenium

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// GetPageSoure returns HTML source of the current page.
func (s *Session) GetPageSoure() string {
	res, err := s.api.executeRequestVoid(
		http.MethodGet, fmt.Sprintf("/session/%s/source", s.id),
	)
	if err != nil {
		handleError(res, err)

		return ""
	}

	if res.Value == nil {
		handleError(nil, errors.New("failed to get page source"))

		return ""
	}

	if v, ok := res.Value.(string); ok {
		return v
	}

	return ""
}

// ExecuteScript executes user defined JavaScript function written as a string.
// ...args are passed to the function as arguments. Function's return value is
// returned as interface{}.
func (s *Session) ExecuteScript(script string, args ...string) interface{} {
	if len(args) == 0 {
		args = []string{}
	}

	script = fmt.Sprintf("return (%s).apply(window, arguments)", script)

	payload := struct {
		Script string   `json:"script"`
		Args   []string `json:"args"`
	}{
		Script: script,
		Args:   args,
	}

	res, err := s.api.executeRequest(
		http.MethodPost, fmt.Sprintf("/session/%s/execute/sync", s.id), payload,
	)
	if err != nil {
		handleError(res, err)

		return s
	}

	return res.Value
}
