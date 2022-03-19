package selenium

import (
	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/client"
)

// Start starts browser driver server and establishes WebDriver session.
func Start(c *client.Client) error {
	err := c.Driver.Launch()
	if err != nil {
		return errors.Wrap(err, "failed to launch driver")
	}

	err = c.StartSession()
	if err != nil {
		return errors.Wrap(err, "failed to start session")
	}

	return nil
}

// TODO: Read selenium.json config.
//nolint: unused,deadcode
func readConfig() {}
