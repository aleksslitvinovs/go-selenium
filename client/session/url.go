package session

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/util"
)

func (s *Session) OpenURL(url string) {
	requestBody := struct {
		URL string `json:"url"`
	}{url}

	res, err := api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/url", s.ID),
		s,
		requestBody,
	)
	if err != nil {
		errRes := res.GetErrorReponse()
		if errRes == nil {
			util.HandleError(s, errors.Wrap(err, "failed to open url"))
		}

		util.HandleResponseError(s, errRes)
	}
}

func (s *Session) GetCurrentURL() string {
	res, err := api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/url", s.ID),
		s,
		struct{}{},
	)
	if err != nil {
		errRes := res.GetErrorReponse()
		if errRes == nil {
			util.HandleError(
				s, errors.Wrap(err, "failed to get current url"),
			)
		}

		util.HandleResponseError(s, errRes)
	}

	return res.Value.(string)
}
