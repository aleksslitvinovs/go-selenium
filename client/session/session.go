package session

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
)

type Session struct {
	URL    string
	Port   int
	ID     string
	Errors []string
}

func (s *Session) GetURL() string {
	return s.URL
}

func (s *Session) GetPort() int {
	return s.Port
}

func NewSession(c api.Requester) (*Session, error) {
	req := struct {
		Capabilities map[string]interface{} `json:"capabilities"`
	}{
		make(map[string]interface{}),
	}

	res, err := api.ExecuteRequestRaw(http.MethodPost, "/session", c, req)

	if err != nil {
		return nil, errors.Wrap(err, "failed to start session")
	}

	var r struct {
		Value struct {
			SessionID    string                 `json:"sessionId"`
			Capabilities map[string]interface{} `json:"capabilities"`
		} `json:"value"`
	}

	if err := json.Unmarshal(res, &r); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response")
	}

	session := &Session{
		URL:  c.GetURL(),
		Port: c.GetPort(),
		ID:   r.Value.SessionID,
	}

	return session, nil

}
func (s *Session) DeleteSession() error {
	_, err := api.ExecuteRequest(
		http.MethodDelete,
		fmt.Sprintf("/session/%s", s.ID),
		s,
		struct{}{},
	)
	if err != nil {
		return errors.Wrap(err, "failed to stop session")
	}

	return nil
}

func (s *Session) Refresh() error {
	_, err := api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/refresh", s.ID),
		s,
		struct{}{},
	)
	if err != nil {
		return errors.Wrap(err, "failed to refresh window")
	}

	return nil
}

func (s *Session) RaiseErrors() string {
	if len(s.Errors) == 0 {
		return ""
	}

	errors := make([]string, 0, len(s.Errors))

	for _, e := range s.Errors {
		errors = append(errors, e)
	}

	return strings.Join(errors, "\n")

}
