package selenium

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/config"
	"github.com/theRealAlpaca/go-selenium/logger"
	"github.com/theRealAlpaca/go-selenium/types"
)

var (
	started = false
)

type client struct {
	api      *api.APIClient
	driver   *driver
	sessions map[*session]bool
}

type Opts struct {
	ConfigPath string
}

// NewClient creates a new client instance with the provided driver. Based on
// the configuration settings, a driver may be started. Optionally, Opts can be
// provided for additional configuration.
func NewClient(d *driver, opts *Opts) (types.Clienter, error) {
	if d == nil {
		return nil, errors.Wrap(
			types.ErrInvalidParameters, "driver cannot be nil",
		)
	}

	if opts == nil {
		opts = &Opts{}
	}

	if !started {
		go gracefulShutdown()

		err := config.ReadConfig(opts.ConfigPath)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read config")
		}

		logger.SetLogLevel(config.Config.LogLevel)

		started = true
	}

	if config.Config.WebDriver.AutoStart {
		err := d.Start(config.Config.WebDriver)
		if err != nil {
			return nil, errors.Wrap(err, "failed to launch driver")
		}
	}

	c := &client{
		api:      &api.APIClient{BaseURL: d.remoteURL},
		driver:   d,
		sessions: make(map[*session]bool),
	}

	return c, nil
}

func gracefulShutdown() {
	var stop = make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	os.Exit(0)
}

func (c *client) MustStop() {
	err := c.Stop()
	if err != nil {
		panic(errors.Wrap(err, "failed to stop driver"))
	}
}

func (c *client) Stop() error {
	var tempErr error

	// Driver must be stopped even if session cannot be deleted.
	defer func() {
		if c.driver == nil {
			return
		}

		err := c.driver.Stop()
		if err != nil {
			tempErr = errors.Wrap(err, "failed to stop driver process")
		}
	}()

	for s, v := range c.sessions {
		if !v {
			continue
		}

		s.DeleteSession()

		if config.Config.RaiseErrorsAutomatically {
			e := s.RaiseErrors()

			if e != "" {
				logger.Errorf("There were issues during execution:\n%s", e)
			}
		}

		delete(c.sessions, s)
	}

	return tempErr
}

func (c *client) RaiseErrors() {
	for s := range c.sessions {
		errors := s.RaiseErrors()

		if len(errors) == 0 {
			continue
		}

		logger.Errorf(
			"Errors occurred in %s session:\n%s\n", s.id, errors,
		)
	}
}

func (c *client) waitUntilIsReady(timeout time.Duration) error {
	endTime := time.Now().Add(timeout)

	for endTime.After(time.Now()) {
		ok, err := c.driver.IsReady(c)
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

	return errors.Errorf(
		"%s timeout exceeded while waiting for driver to be ready",
		timeout.String(),
	)
}
