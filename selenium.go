package selenium

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/config"
	"github.com/theRealAlpaca/go-selenium/driver"
	"github.com/theRealAlpaca/go-selenium/logger"
	"github.com/theRealAlpaca/go-selenium/session"
)

type Opts struct {
	ConfigPath string
}

var started = false

// Start creates a new client instances, starts browser driver server and
// establishes new WebDriver session. Returns a new client instance and the
// established connection.
func Start(opts *Opts) (*client, *session.Session) {
	if opts == nil {
		opts = &Opts{}
	}

	if !started {
		go gracefulShutdown()

		err := config.ReadConfig(opts.ConfigPath)
		if err != nil {
			panic(errors.Wrap(err, "failed to read config"))
		}

		logger.SetLogLevel(config.Config.LogLevel)

		started = true
	}

	d, err := driver.Start(config.Config.WebDriver)
	if err != nil {
		panic(errors.Wrap(err, "failed to launch driver"))
	}

	c := NewClient(d)

	session, err := c.StartNewSession()
	if err != nil {
		panic(errors.Wrap(err, "failed to start session"))
	}

	return c, session
}

func gracefulShutdown() {
	var stop = make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	os.Exit(0)
}
