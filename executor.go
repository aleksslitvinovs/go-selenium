package selenium

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/TylerBrock/colorjson"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/logger"
	"github.com/theRealAlpaca/go-selenium/types"
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

var (
	f = &colorjson.Formatter{
		KeyColor:        color.New(color.FgWhite),
		StringColor:     color.New(color.FgGreen),
		BoolColor:       color.New(color.FgYellow),
		NumberColor:     color.New(color.FgCyan),
		NullColor:       color.New(color.FgMagenta),
		StringMaxLength: 50,
		DisabledColor:   false,
		Indent:          0,
		RawStrings:      true,
	}
)

type APIClient struct {
	BaseURL string
}

func (a *APIClient) ExecuteRequest(
	method, route string, payload interface{},
) (*Response, error) {
	res, reqErr := a.executeRequestRaw(method, route, payload)
	if reqErr != nil {
		if !errors.As(reqErr, &types.ErrFailedRequest) {
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

func (a *APIClient) ExecuteRequestVoid(
	method, route string,
) (*Response, error) {
	return a.ExecuteRequest(method, route, struct{}{})
}

func (a *APIClient) ExecuteRequestCustom(
	method, route string, payload, customResponse interface{},
) error {
	res, err := a.executeRequestRaw(method, route, payload)
	if err != nil {
		return errors.Wrap(err, "failed to execute request")
	}

	if err := json.Unmarshal(res, customResponse); err != nil {
		return errors.Wrap(err, "failed to unmarshal response")
	}

	return nil
}

func (a *APIClient) executeRequestRaw(
	method, route string, payload interface{},
) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal payload")
	}

	url := a.BaseURL + route

	req, err := http.NewRequestWithContext(
		context.Background(), method, url, bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	if Config.LogLevel == logger.DebugLvl {
		logger.Custom(
			color.HiCyanString("-> Request "),
			fmt.Sprintf("%s %s\n\t%s", method, url, formatJSON(body)),
		)
	}

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

	if Config.LogLevel == logger.DebugLvl {
		logger.Custom(
			color.HiGreenString("<- Response "),
			formatJSON(b), "\n\n",
		)
	}

	if getStatusClass(res.StatusCode) != classSuccessful {
		return b, errors.Wrap(types.ErrFailedRequest, response.String())
	}

	return b, nil
}

func formatJSON(body []byte) string {
	var data map[string]interface{}

	//nolint: errcheck
	json.Unmarshal(body, &data)

	b, _ := f.Marshal(data)

	return string(b)
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
