package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/driver"
)

// Describes all possible status code classes.
const (
	ClassInformational = iota + 1
	ClassSuccessful
	ClassRedirection
	ClassClientError
	ClassServerError
)

type errorResponse struct {
	Value errorValue `json:"value"`
}

type errorValue struct {
	Error      string                 `json:"error"`
	Message    string                 `json:"message"`
	Stacktrace string                 `json:"stacktrace"`
	Data       map[string]interface{} `json:"data"`
}

// TODO: Remove driver.Driver from params
func ExecuteRequest(
	method, route string, payload interface{}, d *driver.Driver,
) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to marshal")
	}

	url := fmt.Sprintf("%s:%d%s", d.RemoteURL, d.Port, route)

	req, err := http.NewRequestWithContext(context.Background(), method, url, bytes.NewBuffer(body))
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to create request")
	}

	// TODO: Do optional request logging
	fmt.Printf("Request %q %q'\n", method, url)
	fmt.Printf("Request body: %s\n", body)

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to send request")
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to read response body")
	}
	fmt.Println("Response body: ", string(b))

	if getStatusClass(res.StatusCode) != ClassSuccessful {
		return b, errors.New("failed to execute request")
	}

	return b, nil
}

func getStatusClass(code int) int {
	class := code / 100
	switch class {
	case 1, 2, 3, 4, 5:
		return class
	default:
		return 0
	}
}
