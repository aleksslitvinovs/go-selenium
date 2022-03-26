package session

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	//nolint:godot,gci
	// TODO: fix import cycle
	// "github.com/theRealAlpaca/go-selenium/util"
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
		errRes := res.GetErrorReponse()
		if errRes != nil {
			return "", errors.Wrap(err, "failed to get current url")
		}

		return "", errors.New(res.String())
	}

	return res.Value.(string), nil
}
