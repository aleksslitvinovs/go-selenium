package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/client/session"
	"github.com/theRealAlpaca/go-selenium/driver"
)

type Client struct {
	Driver   *driver.Driver
	Error    error
	Sessions []*session.Session
}

func (c *Client) GetURL() string {
	return c.Driver.RemoteURL
}

func (c *Client) GetPort() int {
	return c.Driver.Port
}

func NewClient(d *driver.Driver) *Client {
	return &Client{
		Driver: d,
	}
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

	for _, s := range c.Sessions {
		err := s.DeleteSession()
		if err != nil {
			return errors.Wrap(err, "failed to stop session")
		}
	}

	return nil
}

func (c *Client) StartSession() (*session.Session, error) {
	if err := waitUntilIsReady(10*time.Second, c); err != nil {
		return &session.Session{}, errors.Wrap(
			err, "driver is not ready to start a new session",
		)
	}

	req := struct {
		Capabilities map[string]interface{} `json:"capabilities"`
	}{
		make(map[string]interface{}),
	}

	res, err := api.ExecuteRequestRaw(http.MethodPost, "/session", c, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start session")
	}

	var r struct {
		Value struct {
			SessionID    string                 `json:"sessionId"`
			Capabilities map[string]interface{} `json:"capabilities"`
		} `json:"value"`
	}

	if err := json.Unmarshal(res, &r); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response")
	}

	session := &session.Session{
		URL: c.Driver.RemoteURL, Port: c.Driver.Port, ID: r.Value.SessionID,
	}

	c.Sessions = append(c.Sessions, session)

	return session, nil
}

func waitUntilIsReady(timeout time.Duration, c *Client) error {
	endTime := time.Now().Add(timeout)

	for {
		if endTime.Before(time.Now()) {
			return errors.New(
				"timeout exceeded while waiting for driver to be ready",
			)
		}

		ok, err := c.IsReady()
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
}
