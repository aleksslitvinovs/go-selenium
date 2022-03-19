package client

import (
	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/driver"
)

type Client struct {
	Driver    *driver.Driver
	Error     error
	SessionID string
}

func NewClientBuilder() *Client {
	return &Client{}
}

func (c *Client) SetDriver(d *driver.Driver) *Client {
	c.Driver = d

	return c
}

func (c *Client) GetDriver() *driver.Driver {
	return c.Driver
}

func (c *Client) Build() *Client {
	return c
}

func (c *Client) Launch() error {
	err := c.Driver.Launch()
	if err != nil {
		return errors.Wrap(err, "failed to launch driver")
	}

	return nil
}

func (c *Client) Stop() error {
	// Driver must be stopped even if session cannot be deleted.
	defer func() error {
		err := c.Driver.Stop()
		if err != nil {
			return errors.Wrap(err, "failed to stop driver")
		}

		return nil
	}() //nolint:errcheck

	err := c.DeleteSession()
	if err != nil {
		return errors.Wrap(err, "failed to stop session")
	}

	return nil
}
