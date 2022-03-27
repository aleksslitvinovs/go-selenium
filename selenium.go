package selenium

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/client"
	"github.com/theRealAlpaca/go-selenium/client/session"
	"github.com/theRealAlpaca/go-selenium/config"
	"github.com/theRealAlpaca/go-selenium/driver"
	"github.com/theRealAlpaca/go-selenium/logger"
)

type Opts struct {
	ConfigPath string
}

// Start starts browser driver server and establishes WebDriver session that
// is returned.
func Start(opts *Opts) *session.Session {
	if opts == nil {
		opts = &Opts{}
	}

	c := client.Client

	go gracefulShutdown()

	err := config.ReadConfig(opts.ConfigPath)
	if err != nil {
		panic(errors.Wrap(err, "failed to read config"))
	}

	logger.SetStringLogLevel(config.Config.LogLevel)

	d, err := driver.Start(config.Config.WebDriver)
	if err != nil {
		panic(errors.Wrap(err, "failed to launch driver"))
	}

	c.Driver = d

	session, err := client.StartNewSession()
	if err != nil {
		panic(errors.Wrap(err, "failed to start session"))
	}

	return session
}

func gracefulShutdown() {
	var stop = make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGINT)

	<-stop

	os.Exit(0)
}
