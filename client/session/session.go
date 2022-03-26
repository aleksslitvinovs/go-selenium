package session

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/util"
)

type Session struct {
	URL        string
	Port       int
	ID         string
	Errors     []string
	KillDriver chan struct{}
}

func (s *Session) GetURL() string {
	return s.URL
}

func (s *Session) GetPort() int {
	return s.Port
}

func (s *Session) AddError(err string) {
	s.Errors = append(s.Errors, err)
}

func NewSession(c api.Requester) (*Session, error) {
	req := struct {
		Capabilities map[string]interface{} `json:"capabilities"`
	}{
		make(map[string]interface{}),
	}

	var r struct {
		Value struct {
			SessionID    string                 `json:"sessionId"`
			Capabilities map[string]interface{} `json:"capabilities"`
		} `json:"value"`
	}

	err := api.ExecuteRequestCustom(http.MethodPost, "/session", c, req, &r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start session")
	}

	session := &Session{
		URL:  c.GetURL(),
		Port: c.GetPort(),
		ID:   r.Value.SessionID,
	}

	return session, nil
}

func (s *Session) Stop() {
	s.DeleteSession() //nolint:errcheck
}

func (s *Session) DeleteSession() {
	res, err := api.ExecuteRequest(
		http.MethodDelete,
		fmt.Sprintf("/session/%s", s.ID),
		s,
		struct{}{},
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			util.HandleResponseError(s, errRes)

			return
		}

		util.HandleError(s, errors.Wrap(err, "failed to delete session"))
	}

	s.KillDriver <- struct{}{}

	return
}

func (s *Session) Refresh() {
	res, err := api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/refresh", s.ID),
		s,
		struct{}{},
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			util.HandleResponseError(s, errRes)

			return
		}

		util.HandleError(s, errors.Wrap(err, "failed to refresh session"))
	}
}

func (s *Session) RaiseErrors() string {
	if len(s.Errors) == 0 {
		return ""
	}

	errors := make([]string, 0, len(s.Errors))

	errors = append(errors, s.Errors...)

	return strings.Join(errors, "\n")
}
