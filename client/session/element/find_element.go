package element

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/client/session"
)

func (e *Element) FindElement(s *session.Session) (string, error) {
	res, err := api.ExecuteRequestRaw(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element", s.ID),
		s,
		e,
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to send request to get element")
	}

	var response struct {
		Value map[string]string `json:"value"`
	}

	if err := json.Unmarshal(res, &response); err != nil {
		return "", errors.Wrap(err, "failed to unmarshal response")
	}

	for _, v := range response.Value {
		if v != "" {
			return v, nil
		}
	}

	return "", errors.New("failed to get element id")
}
