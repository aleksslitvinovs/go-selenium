package selenium

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/client"
	"github.com/theRealAlpaca/go-selenium/client/session"
)

// Start starts browser driver server and establishes WebDriver session that
// is returned.
func Start(c *client.Client) (*session.Session, error) {
	go gracefulShutdown(c)

	err := c.Driver.Launch()
	if err != nil {
		return nil, errors.Wrap(err, "failed to launch driver")
	}

	session, err := c.StartSession()
	if err != nil {
		return nil, errors.Wrap(err, "failed to start session")
	}

	return session, nil
}

func gracefulShutdown(c *client.Client) {
	exit := make(chan os.Signal, 1)

	signal.Notify(exit, syscall.SIGINT)

	<-exit

	if err := c.Stop(); err != nil {
		panic(errors.Wrap(err, "failed to stop client"))
	}

	os.Exit(0)
}

// TODO: Read selenium.json config.
//nolint: unused,deadcode
func readConfig() {}
