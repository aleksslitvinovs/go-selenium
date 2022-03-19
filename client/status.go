package client

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func (c *Client) IsReady() (bool, error) {
	type successResponse struct {
		Ready   bool   `json:"ready"`
		Message string `json:"message"`
	}

	type response struct {
		Value successResponse `json:"value"`
	}

	res, err := ExecuteRequest(http.MethodGet, "/status", struct{}{}, c.Driver)
	if err != nil {
		return false, errors.Wrap(err, "failed to get status")
	}

	var r response

	err = json.Unmarshal(res, &r)
	if err != nil {
		return false, errors.Wrap(err, "failed to unmarshal response")
	}

	return r.Value.Ready, nil
}
