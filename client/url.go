package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

func (c *Client) OpenURL(url string) error {
	requestBody := struct {
		URL string `json:"url"`
	}{
		url,
	}

	_, err := ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/url", c.SessionID),
		requestBody,
		c.Driver,
	)
	if err != nil {
		return errors.Wrap(err, "failed to open url")
	}

	return nil
}

func (c *Client) GetURL() (string, error) {
	res, err := ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/url", c.SessionID),
		struct{}{},
		c.Driver,
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to get url")
	}

	var response struct {
		Value string `json:"value"`
	}

	err = json.Unmarshal(res, &response)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal response")
	}

	return response.Value, nil
}
