package client

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
)

func (c *Client) IsReady() (bool, error) {
	res, err := api.ExecuteRequestRaw(http.MethodGet, "/status", c, struct{}{})
	if err != nil {
		return false, errors.Wrap(err, "failed to get status")
	}

	type response struct {
		Value struct {
			Ready   bool   `json:"ready"`
			Message string `json:"message"`
		} `json:"value"`
	}

	var r response

	if err := json.Unmarshal(res, &r); err != nil {
		return false, errors.Wrap(err, "failed to unmarshal response")
	}

	return r.Value.Ready, nil
}
