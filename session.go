package selenium

import (
	"strings"

	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/types"
)

type session struct {
	id string

	// killDriver is a channel used to kill the client and the driver. Used
	// internally and should not be used by the user.
	killDriver chan struct{}

	url    string
	errors []string
	api    *api.APIClient
}

var _ types.Sessioner = (*session)(nil)

func (s *session) GetID() string {
	return s.id
}
func (s *session) AddError(err string) {
	s.errors = append(s.errors, err)
}

func (s *session) RaiseErrors() string {
	if len(s.errors) == 0 {
		return ""
	}

	errors := make([]string, 0, len(s.errors))

	errors = append(errors, s.errors...)

	return strings.Join(errors, "\n")
}
