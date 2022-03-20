package session

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
)

type Session struct {
	URL  string
	Port int
	ID   string
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
