package selenium

import (
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/logger"
)

type sessionStore struct {
	mu       sync.Mutex
	sessions map[*Session]bool
}

type clientParams struct {
	api    *apiClient
	driver *Driver
	ss     *sessionStore
}

// Opts contains configuration options for the client.
// TODO: Allow overwriting config values.
type Opts struct {
	ConfigDirectory string
}

var client *clientParams

// SetClient creates a new client instance with the provided driver. Based on
// the configuration settings, a driver may be started. Optionally, Opts can be
// provided for additional configuration.
func SetClient(d *Driver, opts *Opts) error {
	if client != nil && client.driver != nil {
		return nil
	}

	if opts == nil {
		opts = &Opts{}
	}

	go gracefulShutdown()

	err := readConfig(opts.ConfigDirectory)
	if err != nil {
		return errors.Wrap(err, "failed to read config")
	}

	logger.SetLogLevel(config.LogLevel)

	if d == nil {
		err := downloadDriver(parseDriver(config.WebDriver.Browser))
		if err != nil {
			return errors.Wrap(err, "failed to download chromedriver")
		}

		d, err = NewDriver(
			config.WebDriver.BinaryPath, config.WebDriver.RemoteURL,
		)
		if err != nil {
			return errors.Wrap(
				err, "failed to create browser default driver",
			)
		}
	}

	if !config.WebDriver.ManualStart {
		err := d.Start(config.WebDriver.Timeout)
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

func parseDriver(driverName string) string {
	switch strings.ToLower(driverName) {
	case "chrome", "chromedriver":
		return chromedriver
	case "firefox", "geckodriver":
		return geckodriver
	default:
		return driverName
	}
}

// MustStopClient is a convenience function that wraps StopClient and panics in
// case an error is encountered.
func MustStopClient() {
	if client == nil {
		return
	}

	err := StopClient()
	if err != nil {
		panic(errors.Wrap(err, "failed to stop driver"))
	}
}

// StopClient stops the client and its driver.
func StopClient() error {
	var tempErr error

	// Driver must be stopped even if session cannot be deleted.
	defer func() {
		if client.driver == nil {
			return
		}

		err := client.driver.stop()
		if err != nil {
			tempErr = errors.Wrap(err, "failed to stop driver process")
		}
	}()

	for s, v := range client.ss.sessions {
		if v {
			s.DeleteSession()
		}

		if !config.RaiseErrorsManually {
			e := s.RaiseErrors()

			if e != "" {
				logger.Errorf("There were issues during execution:\n%s", e)
			}
		}

		client.ss.mu.Lock()
		delete(client.ss.sessions, s)
		client.ss.mu.Unlock()
	}

	return tempErr
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
