package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/client/session"
	"github.com/theRealAlpaca/go-selenium/config"
	"github.com/theRealAlpaca/go-selenium/driver"
)

type client struct {
	Driver   *driver.Driver
	Config   *config.Config
	Error    error
	Sessions map[*session.Session]bool
}

var Client = &client{Sessions: make(map[*session.Session]bool)}

func (c *client) GetURL() string {
	return Client.Driver.RemoteURL
}

func (c *client) GetPort() int {
	return Client.Driver.Port
}

func NewClient(d *driver.Driver) *client {
	Client.Driver = d

	return Client
}

func StartNewSession() (*session.Session, error) {
	if err := waitUntilIsReady(10 * time.Second); err != nil {
		return &session.Session{}, errors.Wrap(
			err, "driver is not ready to start a new session",
		)
	}

	req := struct {
		Capabilities map[string]interface{} `json:"capabilities"`
	}{
		make(map[string]interface{}),
	}

	res, err := api.ExecuteRequestRaw(http.MethodPost, "/session", Client, req)
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

	session := &session.Session{
		Config: Client.Config,
		URL:    Client.Driver.RemoteURL,
		Port:   Client.Driver.Port,
		ID:     r.Value.SessionID,
	}

	Client.Sessions[session] = true

	return session, nil
}

func DeleteSession(s *session.Session) error {
	if err := s.DeleteSession(); err != nil {
		return errors.Wrap(err, "failed to delete session")
	}

	delete(Client.Sessions, s)

	return nil
}

func Stop() error {
	// Driver must be stopped even if session cannot be deleted.
	defer func() error {
		err := Client.Driver.Stop()
		if err != nil {
			return errors.Wrap(err, "failed to stop driver")
		}

		return nil
	}() //nolint:errcheck

	for s, v := range Client.Sessions {
		if !v {
			continue
		}

		err := s.DeleteSession()
		if err != nil {
			return errors.Wrap(err, "failed to stop session")
		}

		if Client.Config.RaiseErrorsAutomaticatically {
			fmt.Println(s.RaiseErrors())
		}

		delete(Client.Sessions, s)
	}

	return nil
}

func (c *client) RaiseErrors() {
	for s := range c.Sessions {
		errors := s.RaiseErrors()

		if len(errors) == 0 {
			continue
		}

		fmt.Printf(
			"Errors occurred in session %s: \n%s\n",
			s.ID, strings.Join(errors, "\n"),
		)
	}
}

func waitUntilIsReady(timeout time.Duration) error {
	endTime := time.Now().Add(timeout)

	for {
		if endTime.Before(time.Now()) {
			return errors.New(
				"timeout exceeded while waiting for driver to be ready",
			)
		}

		ok, err := IsReady()
		if err != nil {
			fmt.Println(err.Error())

			netErr := errors.New("dial tcp")
			if errors.As(err, &netErr) {
				continue
			}

			return errors.Wrap(err, "failed to get status")
		}

		if ok {
			return nil
		}

		time.Sleep(500 * time.Millisecond)
	}
}
