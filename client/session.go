package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

func (c *Client) StartSession() error {
	err := waitUntilIsReady(10*time.Second, c)
	if err != nil {
		return errors.Wrap(err, "driver is not ready to start a new session")
	}

	type successResponse struct {
		SessionID    string                 `json:"sessionId"`
		Capabilities map[string]interface{} `json:"capabilities"`
	}

	type response struct {
		Value successResponse `json:"value"`
	}

	req := struct {
		Capabilities map[string]interface{} `json:"capabilities"`
	}{
		make(map[string]interface{}),
	}

	res, err := ExecuteRequest(http.MethodPost, "/session", req, c.Driver)
	if err != nil {
		return errors.Wrap(err, "failed to start session")
	}

	var r response

	err = json.Unmarshal(res, &r)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal response")
	}

	c.SessionID = r.Value.SessionID

	return nil
}

func (c *Client) DeleteSession() error {
	_, err := ExecuteRequest(
		http.MethodDelete,
		fmt.Sprintf("/session/%s", c.SessionID),
		struct{}{},
		c.Driver,
	)
	if err != nil {
		return errors.Wrap(err, "failed to stop session")
	}

	return nil
}

func (c *Client) Refresh() error {
	_, err := ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/refresh", c.SessionID),
		struct{}{},
		c.Driver,
	)
	if err != nil {
		return errors.Wrap(err, "failed to refresh window")
	}

	return nil
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

		time.Sleep(time.Second)
	}
}
