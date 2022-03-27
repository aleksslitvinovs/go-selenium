package session

import (
	"strings"

	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/util"
)

type Session struct {
	// Navigation is a helper field to access navigation methods.
	Navigation Navigator

	// Context is a helper field to access browsering context related methods.
	Context Contexter

	// KillDriver is a channel used to kill the client and the driver. Used
	// internally and should not be used by the user.
	KillDriver chan struct{}

	url    string
	port   int
	ID     string
	errors []string
}

var (
	_ Navigator      = (*Session)(nil)
	_ Contexter      = (*Session)(nil)
	_ util.Sessioner = (*Session)(nil)
	_ api.Requester  = (*Session)(nil)
)

func NewSession(c api.Requester) (*Session, error) {
	id := CreateSession(c)

	session := &Session{
		url:  c.GetURL(),
		port: c.GetPort(),
		ID:   id,
	}

	session.Navigation = session
	session.Context = session

	return session, nil
}

func (s *Session) GetURL() string {
	return s.url
}

func (s *Session) GetPort() int {
	return s.port
}

func (s *Session) GetErrors() []string {
	return s.errors
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
