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
		if errRes := res.GetErrorReponse(); errRes != nil {
			util.HandleResponseError(s, errRes)

			return
		}

		util.HandleError(s, errors.Wrap(err, "failed to open URL"))
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
		if errRes := res.GetErrorReponse(); errRes != nil {
			util.HandleResponseError(s, errRes)

			return ""
		}

		util.HandleError(s, errors.Wrap(err, "failed to get current URL"))
	}

	return res.Value.(string)
}

func (s *Session) Refresh() {
	res, err := api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/refresh", s.ID),
		s,
		struct{}{},
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			util.HandleResponseError(s, errRes)

			return
		}

		util.HandleError(s, errors.Wrap(err, "failed to refresh session"))
	}
}

func (s *Session) Back() {
	res, err := api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/back", s.ID),
		s,
		struct{}{},
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			util.HandleResponseError(s, errRes)

			return
		}

		util.HandleError(s, errors.Wrap(err, "failed to go back"))
	}
}

func (s *Session) Forward() {
	res, err := api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/forward", s.ID),
		s,
		struct{}{},
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			util.HandleResponseError(s, errRes)

			return
		}

		util.HandleError(s, errors.Wrap(err, "failed to go forward"))
	}
}

func (s *Session) GetTitle() string {
	res, err := api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/title", s.ID),
		s,
		struct{}{},
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			util.HandleResponseError(s, errRes)

			return ""
		}

		util.HandleError(s, errors.Wrap(err, "failed to get page title"))
	}

	return res.Value.(string)
}
