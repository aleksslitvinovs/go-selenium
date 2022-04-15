package selenium

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/logger"
)

type client struct {
	api      *APIClient
	driver   *Driver
	sessions map[*Session]bool
}

type Opts struct {
	ConfigPath string
}

var Client *client

// NewClient creates a new client instance with the provided driver. Based on
// the configuration settings, a driver may be started. Optionally, Opts can be
// provided for additional configuration.
func StartClient(d *Driver, opts *Opts) (*client, error) {
	if Client != nil {
		return Client, nil
	}

	if opts == nil {
		opts = &Opts{}
	}

	wg := &sync.WaitGroup{}

	if d == nil {
		wg.Add(1)

		err := downloadDriver(wg, chromedriver)
		if err != nil {
			return nil, errors.Wrap(err, "failed to download chromedriver")
		}

		d, err = NewDriver(chromedriver, "http://localhost:4445")
		if err != nil {
			return nil, errors.Wrap(
				err, "failed to create browser default driver",
			)
		}
	}

	go gracefulShutdown()

	err := ReadConfig(opts.ConfigPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config")
	}

	logger.SetLogLevel(Config.LogLevel)

	wg.Wait()

	if !Config.WebDriver.ManualStart {
		err := d.Start(Config.WebDriver)
		if err != nil {
			return nil, errors.Wrap(err, "failed to launch driver")
		}
	}

	Client = &client{
		api:      &APIClient{BaseURL: d.remoteURL},
		driver:   d,
		sessions: make(map[*Session]bool),
	}

	return Client, nil
}

func gracefulShutdown() {
	var stop = make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	os.Exit(0)
}

func MustStopClient() {
	if Client == nil {
		return
	}

	err := StopClient()
	if err != nil {
		panic(errors.Wrap(err, "failed to stop driver"))
	}
}

func StopClient() error {
	var tempErr error

	// Driver must be stopped even if session cannot be deleted.
	defer func() {
		if Client.driver == nil {
			return
		}

		err := Client.driver.Stop()
		if err != nil {
			tempErr = errors.Wrap(err, "failed to stop driver process")
		}
	}()

	for s, v := range Client.sessions {
		if v {
			// TODO: Handle already deleted session
			s.DeleteSession()
		}

		if Config.RaiseErrorsAutomatically {
			e := s.RaiseErrors()

			if e != "" {
				logger.Errorf("There were issues during execution:\n%s", e)
			}
		}

		delete(Client.sessions, s)
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
			"Errors occurred in %s session:\n%s\n", s.GetID(), errors,
		)
	}
}

func (c *client) waitUntilIsReady(timeout time.Duration) error {
	endTime := time.Now().Add(timeout)

	for endTime.After(time.Now()) {
		ok, err := c.driver.IsReady(c)
		if err != nil {
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
