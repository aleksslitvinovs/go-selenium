package selenium

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/client"
	"github.com/theRealAlpaca/go-selenium/client/session"
	"github.com/theRealAlpaca/go-selenium/config"
)

// Start starts browser driver server and establishes WebDriver session that
// is returned.
func Start() (*session.Session, error) {
	c := client.Client

	go gracefulShutdown()

	conf, err := config.ReadConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config")
	}

	c.Config = conf

	err = c.Driver.Launch()
	if err != nil {
		return nil, errors.Wrap(err, "failed to launch driver")
	}

	session, err := client.StartNewSession()
	if err != nil {
		return nil, errors.Wrap(err, "failed to start session")
	}

	return session, nil
}

func gracefulShutdown() {
	var stop = make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGINT)

	<-stop

	if err := client.Stop(); err != nil {
		panic(errors.Wrap(err, "failed to stop client"))
	}

	os.Exit(0)
}

// TODO: Read selenium.json config.
//nolint: unused,deadcode
func readConfig() {}
