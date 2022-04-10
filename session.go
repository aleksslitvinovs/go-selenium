package selenium

import (
	"strings"

	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/types"
)

type Session struct {
	id string

	// killDriver is a channel used to kill the client and the driver. Used
	// internally and should not be used by the user.
	killDriver chan struct{}

	url string
	// TODO: Maybe create a custom struct for handling error types. Maybe just
	// an alias to string? Maybe could implement Error interface?
	errors []string
	api    *api.APIClient
}

var _ types.Sessioner = (*Session)(nil)

func (s *Session) GetID() string {
	return s.id
}
func (s *Session) AddError(err string) {
	s.errors = append(s.errors, err)
}

func (s *Session) RaiseErrors() string {
	if len(s.errors) == 0 {
		return ""
	}

	errors := make([]string, 0, len(s.errors))

	errors = append(errors, s.errors...)

	return strings.Join(errors, "\n")
}
