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

	err := json.Unmarshal(data, &res)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal response")
	}

	if _, ok := res.Value.(map[string]string); ok {
		var response struct {
			Value ExpandedResponse `json:"value"`
		}

		if err := json.Unmarshal(data, &response); err != nil {
			return errors.Wrap(err, "failed to unmarshal response")
		}

		r.Value = response.Value

		return nil
	}

	if values, ok := res.Value.(map[string]interface{}); ok {
		for k, v := range values {
			if strings.HasPrefix(k, "element") {
				r.Value = map[string]string{k: v.(string)}

				return nil
			}
		}

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
