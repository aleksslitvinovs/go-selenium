package client

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/client/session"
	"github.com/theRealAlpaca/go-selenium/config"
	"github.com/theRealAlpaca/go-selenium/driver"
	"github.com/theRealAlpaca/go-selenium/logger"
)

type client struct {
	Driver   *driver.Driver
	Error    error
	Sessions map[*session.Session]bool
}

var (
	Client = &client{Sessions: make(map[*session.Session]bool)}
	done   = make(chan struct{})
)

func (c *client) GetURL() string {
	return Client.Driver.RemoteURL
}

func (c *client) GetPort() int {
	return Client.Driver.Port
}

func StartNewSession() (*session.Session, error) {
	if err := waitUntilIsReady(10 * time.Second); err != nil {
		return &session.Session{}, errors.Wrap(
			err, "driver is not ready to start a new session",
		)
	}

	s, err := session.NewSession(Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start new session")
	}

	s.KillDriver = done

	Client.Sessions[s] = true

	go sessionListener(s)

	return s, nil
}

func sessionListener(s *session.Session) {
	if len(Client.Sessions) == 0 {
		Stop()
	}

	<-s.KillDriver

	delete(Client.Sessions, s)

	if len(Client.Sessions) == 0 {
		Stop()
	}
}

func DeleteSession(s *session.Session) error {
	s.DeleteSession()

	delete(Client.Sessions, s)

	return nil
}

func Stop() {
	exitCode := 0

	// Driver must be stopped even if session cannot be deleted.
	defer func() {
		err := Client.Driver.Stop()
		if err != nil {
			panic(errors.Wrap(err, "failed to stop driver"))
		}

		os.Exit(exitCode)
	}()

	for s, v := range Client.Sessions {
		if !v {
			continue
		}

		s.DeleteSession()

		if len(s.GetErrors()) != 0 {
			exitCode = 1
		}

		if config.Config.RaiseErrorsAutomatically {
			e := s.RaiseErrors()

			if e != "" {
				logger.Errorf("There were issues during execution:\n%s", e)

				exitCode = 1
			}
		}

		delete(Client.Sessions, s)
	}
}

func (c *client) RaiseErrors() {
	for s := range c.Sessions {
		errors := s.RaiseErrors()

		if len(errors) == 0 {
			continue
		}

		fmt.Printf(
			"Errors occurred in %s session:\n%s\n", s.ID, errors,
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

		ok, err := driver.IsReady(Client)
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
