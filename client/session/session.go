package session

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/config"
)

type Session struct {
	Config *config.Config
	URL    string
	Port   int
	ID     string
	Errors []error
}

func (s *Session) GetURL() string {
	return s.URL
}

func (s *Session) GetPort() int {
	return s.Port
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

func (s *Session) RaiseErrors() []string {
	if len(s.Errors) == 0 {
		return []string{}
	}

	errors := make([]string, 0, len(s.Errors))

	for _, err := range s.Errors {
		errors = append(errors, err.Error())
	}

	return errors
}
