package selenium

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/types"
)

type handleType string

var (
	tab    handleType = "tab"
	window handleType = "window"
)

// OpenURL opens a new window with the given URL.
func (s *Session) OpenURL(url string) {
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
			HandleResponseError(errRes)

			return
		}

		HandleError(errors.Wrap(err, "failed to open URL"))
	}
}

// GetCurrentURL returns the current URL of the browsing context.
func (s *Session) GetCurrentURL() string {
	res, err := s.api.ExecuteRequestVoid(
		http.MethodGet,
		fmt.Sprintf("/session/%s/url", s.id),
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			HandleResponseError(errRes)

			return ""
		}

		HandleError(errors.Wrap(err, "failed to get current URL"))
	}

	return res.Value.(string)
}

// Refresh refreshes the current page.
func (s *Session) Refresh() {
	res, err := s.api.ExecuteRequestVoid(
		http.MethodPost,
		fmt.Sprintf("/session/%s/refresh", s.id),
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			HandleResponseError(errRes)

			return
		}

		HandleError(errors.Wrap(err, "failed to refresh session"))
	}
}

// Back navigates back in the browser history.
func (s *Session) Back() {
	res, err := s.api.ExecuteRequestVoid(
		http.MethodPost,
		fmt.Sprintf("/session/%s/back", s.id),
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			HandleResponseError(errRes)

			return
		}

		HandleError(errors.Wrap(err, "failed to go back"))
	}
}

// Forward navigates forward in the browser history.
func (s *Session) Forward() {
	res, err := s.api.ExecuteRequestVoid(
		http.MethodPost,
		fmt.Sprintf("/session/%s/forward", s.id),
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			HandleResponseError(errRes)

			return
		}

		HandleError(errors.Wrap(err, "failed to go forward"))
	}
}

// GetTitle returns the current page title.
func (s *Session) GetTitle() string {
	res, err := s.api.ExecuteRequestVoid(
		http.MethodGet,
		fmt.Sprintf("/session/%s/title", s.id),
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			HandleResponseError(errRes)

			return ""
		}

		HandleError(errors.Wrap(err, "failed to get page title"))
	}

	return res.Value.(string)
}

func (s *Session) GetWindowHandle() string {
	res, err := s.api.ExecuteRequestVoid(
		http.MethodGet,
		fmt.Sprintf("/session/%s/window", s.id),
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			HandleResponseError(errRes)

			return ""
		}

		HandleError(errors.Wrap(err, "failed to get window handle"))
	}

	return res.Value.(string)
}

func (s *Session) GetWindowHandles() []string {
	res, err := s.api.ExecuteRequestVoid(
		http.MethodGet,
		fmt.Sprintf("/session/%s/window/handles", s.id),
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			HandleResponseError(errRes)

			return []string{}
		}

		HandleError(errors.Wrap(err, "failed to get window handles"))
	}

	result := res.Value.([]interface{})
	handles := make([]string, 0, len(result))

	for _, v := range result {
		handles = append(handles, v.(string))
	}

	return handles
}

// If there are no open browsing contexts left, the session is closed.
func (s *Session) CloseWindow() {
	res, err := s.api.ExecuteRequestVoid(
		http.MethodDelete,
		fmt.Sprintf("/session/%s/window", s.id),
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			HandleResponseError(errRes)

			return
		}

		HandleError(errors.Wrap(err, "failed to close window"))
	}
}

func (s *Session) SwitchHandle(handle string) {
	payload := struct {
		Handle string `json:"handle"`
	}{handle}

	res, err := s.api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/window", s.id),
		payload,
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			HandleResponseError(errRes)

			return
		}

		HandleError(errors.Wrap(err, "failed to switch to handle"))
	}
}

func (s *Session) NewTab() *types.Handle {
	return s.newWindowWithType(tab)
}

func (s *Session) NewWindow() *types.Handle {
	return s.newWindowWithType(window)
}

func (s *Session) newWindowWithType(ht handleType) *types.Handle {
	payload := struct {
		HandleType string `json:"type"`
	}{string(ht)}

	var response struct {
		Value types.Handle `json:"value"`
	}

	err := s.api.ExecuteRequestCustom(
		http.MethodPost,
		fmt.Sprintf("/session/%s/window/new", s.id),
		payload,
		&response,
	)
	if err != nil {
		HandleError(
			errors.Wrapf(err, "failed to open new %s", string(ht)),
		)
	}

	return &response.Value
}

// FIXME: Allow passing handle ID instead of int
// https://www.w3.org/TR/webdriver/#switch-to-frame
func (s *Session) SwitchToFrame(id int) {
	payload := struct {
		ID int `json:"id"`
	}{id}

	res, err := s.api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/frame", s.id),
		payload,
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			HandleResponseError(errRes)

			return
		}

		HandleError(errors.Wrap(err, "failed to get window handles"))
	}
}

func (s *Session) SwitchToParentFrame() {
	res, err := s.api.ExecuteRequestVoid(
		http.MethodPost,
		fmt.Sprintf("/session/%s/frame/parent", s.id),
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			HandleResponseError(errRes)

			return
		}

		HandleError(errors.Wrap(err, "failed to get window handles"))
	}
}
