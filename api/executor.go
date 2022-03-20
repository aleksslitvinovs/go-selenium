package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
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

type Requester interface {
	GetURL() string
	GetPort() int
}

type Response struct {
	Value interface{} `json:"value"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func (r *Response) String() string {
	switch v := r.Value.(type) {
	case ErrorResponse:
		return fmt.Sprintf(
			`"error": %q, "message": %q`,
			v.Error, v.Message,
		)
	default:
		return fmt.Sprintf(`"value": %q`, r.Value)
	}
}

func (r *Response) UnmarshalJSON(data []byte) error {
	var res struct {
		Value interface{} `json:"value"`
	}

	err := json.Unmarshal(data, &res)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal response")
	}

	if _, ok := res.Value.(map[string]interface{}); ok {
		var errResponse struct {
			Value ErrorResponse `json:"value"`
		}

		err = json.Unmarshal(data, &errResponse)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal error response")
		}

		r.Value = errResponse.Value

		return nil
	}

	r.Value = res.Value

	return nil
}

func ExecuteRequestRaw(
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

	// TODO: Do optional request logging
	fmt.Printf("Request %q %q'\n", method, url)
	fmt.Printf("Request body: %s\n", body)

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	var response Response

	err = json.Unmarshal(b, &response)
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to unmarshal response")
	}

	fmt.Println("Raw response", string(b))

	if getStatusClass(res.StatusCode) != classSuccessful {
		return nil, errors.Errorf(
			"failed to execute request, response body {%s}", response.String(),
		)
	}

	return b, nil
}

func ExecuteRequest(
	method, route string, r Requester, payload interface{},
) (*Response, error) {
	res, err := ExecuteRequestRaw(method, route, r, payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute request")
	}

	var response Response

	err = json.Unmarshal(res, &response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response")
	}

	return &response, nil
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
