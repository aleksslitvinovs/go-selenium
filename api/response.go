package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type Response struct {
	Value interface{} `json:"value"`
}

type ExpandedResponse struct {
	Elements map[string]string `json:"-"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func (r *Response) GetValue() interface{} {
	return r.Value
}

func (r *Response) GetErrorReponse() *ErrorResponse {
	if r == nil {
		return nil
	}

	if r.Value == nil {
		return nil
	}

	if errRes, ok := r.Value.(ErrorResponse); ok {
		return &errRes
	}

	return nil
}

func (errRes *ErrorResponse) String() string {
	return fmt.Sprintf(
		`"error": %q, "message": %q`, errRes.Error, errRes.Message,
	)
}

func (r *Response) String() string {
	switch v := r.Value.(type) {
	case ErrorResponse:
		return v.String()
	default:
		return fmt.Sprintf(`"value": %q`, r.Value)
	}
}

func (r *Response) UnmarshalJSON(data []byte) error {
	var res struct {
		Value interface{} `json:"value"`
	}

	if err := json.Unmarshal(data, &res); err != nil {
		return errors.Wrap(err, "failed to unmarshal response")
	}

	switch values := res.Value.(type) {
	case map[string]interface{}:
		for k, v := range values {
			if strings.HasPrefix(k, "element") {
				r.Value = map[string]string{k: v.(string)}

				return nil
			}
		}

		var errResponse struct {
			Value ErrorResponse `json:"value"`
		}

		if err := json.Unmarshal(data, &errResponse); err != nil {
			return errors.Wrap(err, "failed to unmarshal error response")
		}

		r.Value = errResponse.Value

		return nil
	case []interface{}:
		var response struct {
			Value []interface{} `json:"value"`
		}

		if err := json.Unmarshal(data, &response); err != nil {
			return errors.Wrap(err, "failed to unmarshal response")
		}

		r.Value = response.Value
	default:
		r.Value = res.Value
	}

	return nil
}
