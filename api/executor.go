package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/logger"
)

// Describes possible response status code classes.
//nolint:deadcode,varcheck
const (
	classInformational = iota + 1
	classSuccessful
	classRedirection
	classClientError
	classServerError
)

var ErrFailedRequest = errors.New("failed to execute request")

type Requester interface {
	GetURL() string
	GetPort() int
}

func ExecuteRequest(
	method, route string, r Requester, payload interface{},
) (*Response, error) {
	res, reqErr := executeRequestRaw(method, route, r, payload)
	if reqErr != nil {
		if !errors.As(reqErr, &ErrFailedRequest) {
			return nil, errors.Wrap(reqErr, "failed to execute request")
		}
	}

	var response Response

	err := json.Unmarshal(res, &response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response")
	}

	if reqErr != nil {
		return &response, reqErr
	}

	return &response, nil
}

func ExecuteRequestVoid(method, route string, r Requester) (*Response, error) {
	return ExecuteRequest(method, route, r, struct{}{})
}

func ExecuteRequestCustom(
	method, route string, r Requester, payload, customResponse interface{},
) error {
	res, err := executeRequestRaw(method, route, r, payload)
	if err != nil {
		return errors.Wrap(err, "failed to execute request")
	}

	if err := json.Unmarshal(res, customResponse); err != nil {
		return errors.Wrap(err, "failed to unmarshal response")
	}

	return nil
}

func executeRequestRaw(
	method, route string, r Requester, payload interface{},
) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal payload")
	}

	url := fmt.Sprintf("%s:%d%s", r.GetURL(), r.GetPort(), route)

	req, err := http.NewRequestWithContext(
		context.Background(), method, url, bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	logger.Debugf("Request: %q '%q'\n\t%s", method, url, body)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	var response Response

	// Unmarshaling is needed to format errors
	err = json.Unmarshal(b, &response)
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to unmarshal response")
	}

	if getStatusClass(res.StatusCode) != classSuccessful {
		return b, errors.Wrap(ErrFailedRequest, response.String())
	}

	logger.Debugf("Response: %s", string(b))

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
