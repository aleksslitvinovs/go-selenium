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

type sessionStore struct {
	sessions map[*Session]bool
	mu       sync.Mutex
}
type clientParams struct {
	api    *apiClient
	driver *Driver
	ss     *sessionStore
}

type Opts struct {
	ConfigPath string
}

var client *clientParams

// NewClient creates a new client instance with the provided driver. Based on
// the configuration settings, a driver may be started. Optionally, Opts can be
// provided for additional configuration.
func SetClient(d *Driver, opts *Opts) error {
	if client != nil {
		return nil
	}

	if opts == nil {
		opts = &Opts{}
	}

	go gracefulShutdown()

	err := readConfig(opts.ConfigPath)
	if err != nil {
		return errors.Wrap(err, "failed to read config")
	}

	logger.SetLogLevel(config.LogLevel)

	wg := &sync.WaitGroup{}

	if d == nil {
		wg.Add(1)

		err := downloadDriver(wg, chromedriver)
		if err != nil {
			return errors.Wrap(err, "failed to download chromedriver")
		}

		d, err = NewDriver(chromedriver, "http://localhost:4445")
		if err != nil {
			return errors.Wrap(
				err, "failed to create browser default driver",
			)
		}
	}

	wg.Wait()

	if !config.WebDriver.ManualStart {
		err := d.Start(config.WebDriver)
		if err != nil {
			return errors.Wrap(err, "failed to launch driver")
		}
	}

	client = &clientParams{
		api:    &apiClient{baseURL: d.remoteURL},
		driver: d,
		ss:     &sessionStore{sessions: make(map[*Session]bool)},
	}

	return nil
}

func gracefulShutdown() {
	var stop = make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	os.Exit(0)
}

func MustStopClient() {
	if client == nil {
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
		if client.driver == nil {
			return
		}

		err := client.driver.Stop()
		if err != nil {
			tempErr = errors.Wrap(err, "failed to stop driver process")
		}
	}()

	client.ss.mu.Lock()
	defer client.ss.mu.Unlock()

	for s, v := range client.ss.sessions {
		if v {
			// TODO: Handle already deleted session
			s.DeleteSession()
		}

		if config.RaiseErrorsAutomatically {
			e := s.RaiseErrors()

			if e != "" {
				logger.Errorf("There were issues during execution:\n%s", e)
			}
		}

		delete(client.ss.sessions, s)
	}

	return tempErr
}

func (c *clientParams) RaiseErrors() {
	c.ss.mu.Lock()
	defer c.ss.mu.Unlock()

	for s := range c.ss.sessions {
		errors := s.RaiseErrors()

		if len(errors) == 0 {
			continue
		}

		logger.Errorf(
			"Errors occurred in %s session:\n%s\n", s.GetID(), errors,
		)
	}
}

func (c *clientParams) waitUntilIsReady(timeout time.Duration) error {
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
