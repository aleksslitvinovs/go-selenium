package session

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/util"
)

type Navigatorer interface {
	// OpenURL opens a new window with the given URL.
	OpenURL(url string)

	// GetCurrentURL returns the current URL of the browsing context.
	GetCurrentURL() string

	// Refresh refreshes the current page.
	Refresh()

	// Back navigates back in the browser history.
	Back()

	// Forward navigates forward in the browser history.
	Forward()

	// GetTitle returns the current page title.
	GetTitle() string
}

// OpenURL opens a new window with the given URL.
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

// GetCurrentURL returns the current URL of the browsing context.
func (s *Session) GetCurrentURL() string {
	res, err := api.ExecuteRequestVoid(
		http.MethodGet,
		fmt.Sprintf("/session/%s/url", s.ID),
		s,
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
func (s *Session) Refresh() {
	res, err := api.ExecuteRequestVoid(
		http.MethodPost,
		fmt.Sprintf("/session/%s/refresh", s.ID),
		s,
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
func (s *Session) Back() {
	res, err := api.ExecuteRequestVoid(
		http.MethodPost,
		fmt.Sprintf("/session/%s/back", s.ID),
		s,
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
func (s *Session) Forward() {
	res, err := api.ExecuteRequestVoid(
		http.MethodPost,
		fmt.Sprintf("/session/%s/forward", s.ID),
		s,
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
func (s *Session) GetTitle() string {
	res, err := api.ExecuteRequestVoid(
		http.MethodGet,
		fmt.Sprintf("/session/%s/title", s.ID),
		s,
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
