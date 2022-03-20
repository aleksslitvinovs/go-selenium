package session

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
)

func (s *Session) OpenURL(url string) error {
	requestBody := struct {
		URL string `json:"url"`
	}{url}

	_, err := api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/url", s.ID),
		s,
		requestBody,
	)
	if err != nil {
		return errors.Wrap(err, "failed to open url")
	}

	return nil
}

func (s *Session) GetCurrentURL() (string, error) {
	res, err := api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/url", s.ID),
		s,
		struct{}{},
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to get url")
	}

	return res.Value.(string), nil
}
