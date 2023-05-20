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
	"github.com/aleksslitvinovs/go-selenium/logger"
	"github.com/aleksslitvinovs/go-selenium/types"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

// Describes possible response status code classes.
//
//nolint:deadcode,varcheck
const (
	classInformational = iota + 1
	classSuccessful
	classRedirection
	classClientError
	classServerError
)

var f = &colorjson.Formatter{
	KeyColor:        color.New(color.FgWhite),
	StringColor:     color.New(color.FgGreen),
	BoolColor:       color.New(color.FgYellow),
	NumberColor:     color.New(color.FgCyan),
	NullColor:       color.New(color.FgMagenta),
	StringMaxLength: 50,
	DisabledColor:   false,
	Indent:          0,
	RawStrings:      false,
}

type apiClient struct {
	baseURL string
}

type response struct {
	Value interface{} `json:"value"`
}

//nolint:errname
type errorResponse struct {
	Err     string `json:"error"`
	Message string `json:"message"`
}

func (a *apiClient) executeRequest(
	method, route string, payload interface{},
) (*response, error) {
	res, reqErr := a.executeRequestRaw(method, route, payload)
	if reqErr != nil {
		if !errors.As(reqErr, &types.ErrFailedRequest) {
			return nil, errors.Wrap(reqErr, "failed to execute request")
		}
	}

	var r response

	err := json.Unmarshal(res, &r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response")
	}

	if reqErr != nil {
		return &r, reqErr
	}

	return &r, nil
}

func (a *apiClient) executeRequestVoid(
	method, route string,
) (*response, error) {
	return a.executeRequest(method, route, struct{}{})
}

func (a *apiClient) executeRequestCustom(
	method, route string, payload, customResponse interface{},
) (*response, error) {
	res, err := a.executeRequestRaw(method, route, payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute request")
	}

	err = json.Unmarshal(res, customResponse)
	if err == nil {
		return &response{Value: customResponse}, nil
	}

	var errRes *response

	err = json.Unmarshal(res, &errRes)
	if err != nil {
		return nil, errors.Wrap(
			err, "failed to unmarshal response into errorResponse",
		)
	}

	return errRes, errors.Wrap(err, "failed to unmarshal response")
}

func (a *apiClient) executeRequestRaw(
	method, route string, payload interface{},
) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal payload")
	}

	url := a.baseURL + route

	req, err := http.NewRequestWithContext(
		context.Background(), method, url, bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	if config.LogLevel == logger.DebugLvl {
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

	var r response

	// Unmarshaling is needed to format errors
	err = json.Unmarshal(b, &r)
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to unmarshal response")
	}

	if config.LogLevel == logger.DebugLvl {
		logger.Custom(
			color.HiGreenString("<- Response "),
			formatJSON(b), "\n\n",
		)
	}

	if getStatusClass(res.StatusCode) != classSuccessful {
		return b, errors.Wrap(types.ErrFailedRequest, r.String())
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

func (r *response) String() string {
	switch v := r.Value.(type) {
	case errorResponse:
		return v.String()
	default:
		return fmt.Sprintf(`"value": %q`, r.Value)
	}
}

func (r *response) UnmarshalJSON(data []byte) error {
	var res struct {
		Value interface{} `json:"value"`
	}

	if err := json.Unmarshal(data, &res); err != nil {
		return errors.Wrap(err, "failed to unmarshal response")
	}

	switch values := res.Value.(type) {
	case map[string]interface{}:
		for k, value := range values {
			if strings.HasPrefix(strings.ToLower(k), "element") {
				if v, ok := value.(string); ok {
					r.Value = map[string]string{k: v}

					return nil
				}

				return errors.New("could not convert element to string")
			}
		}

		var errResponse struct {
			Value errorResponse `json:"value"`
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

func (r *response) getErrorReponse() *errorResponse {
	if r == nil {
		return nil
	}

	if r.Value == nil {
		return nil
	}

	if errRes, ok := r.Value.(errorResponse); ok {
		return &errRes
	}

	return nil
}

func (errRes *errorResponse) String() string {
	return fmt.Sprintf(
		`"error": %q, "message": %q`, errRes.Err, errRes.Message,
	)
}

func (errRes *errorResponse) Error() string {
	return errRes.Err
}
