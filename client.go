package selenium

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/config"
	"github.com/theRealAlpaca/go-selenium/driver"
	"github.com/theRealAlpaca/go-selenium/logger"
	"github.com/theRealAlpaca/go-selenium/session"
)

type client struct {
	Driver   *driver.Driver
	Sessions map[*session.Session]bool
}

var (
	done = make(chan struct{})
)

func NewClient(d *driver.Driver) *client {
	if d == nil {
		panic("driver cannot be nil")
	}

	return &client{
		Driver:   d,
		Sessions: make(map[*session.Session]bool),
	}
}

func (c *client) GetURL() string {
	return c.Driver.RemoteURL
}

func (c *client) GetPort() int {
	return c.Driver.Port
}

func (c *client) StartNewSession() (*session.Session, error) {
	if err := c.waitUntilIsReady(10 * time.Second); err != nil {
		return &session.Session{}, errors.Wrap(
			err, "driver is not ready to start a new session",
		)
	}

	s, err := session.NewSession(c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start new session")
	}

	s.KillDriver = done

	c.Sessions[s] = true

	go c.sessionListener(s)

	return s, nil
}

func (c *client) sessionListener(s *session.Session) {
	if len(c.Sessions) == 0 {
		c.Stop()
	}

	<-s.KillDriver

	delete(c.Sessions, s)

	if len(c.Sessions) == 0 {
		c.Stop()
	}
}

func (c *client) Stop() {
	exitCode := 0

	// Driver must be stopped even if session cannot be deleted.
	defer func() {
		err := c.Driver.Stop()
		if err != nil {
			panic(errors.Wrap(err, "failed to stop driver"))
		}

		os.Exit(exitCode)
	}()

	for s, v := range c.Sessions {
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

		delete(c.Sessions, s)
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

func (c *client) waitUntilIsReady(timeout time.Duration) error {
	endTime := time.Now().Add(timeout)

	for {
		if endTime.Before(time.Now()) {
			return errors.New(
				"timeout exceeded while waiting for driver to be ready",
			)
		}

		ok, err := driver.IsReady(c)
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
