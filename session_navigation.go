package selenium

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/util"
)

// OpenURL opens a new window with the given URL.
func (s *session) OpenURL(url string) {
	requestBody := struct {
		URL string `json:"url"`
	}{url}

	res, err := s.api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/url", s.id),
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

// GetCurrentURL returns the current URL of the browsing context.
func (s *session) GetCurrentURL() string {
	res, err := s.api.ExecuteRequestVoid(
		http.MethodGet,
		fmt.Sprintf("/session/%s/url", s.id),
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

// Refresh refreshes the current page.
func (s *session) Refresh() {
	res, err := s.api.ExecuteRequestVoid(
		http.MethodPost,
		fmt.Sprintf("/session/%s/refresh", s.id),
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			util.HandleResponseError(s, errRes)

			return
		}

		util.HandleError(s, errors.Wrap(err, "failed to refresh session"))
	}
}

// Back navigates back in the browser history.
func (s *session) Back() {
	res, err := s.api.ExecuteRequestVoid(
		http.MethodPost,
		fmt.Sprintf("/session/%s/back", s.id),
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			util.HandleResponseError(s, errRes)

			return
		}

		util.HandleError(s, errors.Wrap(err, "failed to go back"))
	}
}

// Forward navigates forward in the browser history.
func (s *session) Forward() {
	res, err := s.api.ExecuteRequestVoid(
		http.MethodPost,
		fmt.Sprintf("/session/%s/forward", s.id),
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			util.HandleResponseError(s, errRes)

			return
		}

		util.HandleError(s, errors.Wrap(err, "failed to go forward"))
	}
}

// GetTitle returns the current page title.
func (s *session) GetTitle() string {
	res, err := s.api.ExecuteRequestVoid(
		http.MethodGet,
		fmt.Sprintf("/session/%s/title", s.id),
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
