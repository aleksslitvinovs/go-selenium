package selenium

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type handleType string

type Handle struct {
	ID   string `json:"handle"`
	Type string `json:"type"`
}

var (
	tab    handleType = "tab"
	window handleType = "window"
)

// OpenURL opens a new window with the given URL.
func (s *Session) OpenURL(url string) *Session {
	requestBody := struct {
		URL string `json:"url"`
	}{url}

	res, err := s.api.executeRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/url", s.id),
		requestBody,
	)
	if err != nil {
		handleError(res, err)
	}

	return s
}

// GetCurrentURL returns the current URL of the browsing context.
func (s *Session) GetCurrentURL() string {
	res, err := s.api.executeRequestVoid(
		http.MethodGet,
		fmt.Sprintf("/session/%s/url", s.id),
	)
	if err != nil {
		handleError(res, err)

		return ""
	}

	return res.Value.(string)
}

// Refresh refreshes the current page.
func (s *Session) Refresh() *Session {
	res, err := s.api.executeRequestVoid(
		http.MethodPost,
		fmt.Sprintf("/session/%s/refresh", s.id),
	)
	if err != nil {
		handleError(res, err)
	}

	return s
}

// Back navigates back in the browser history.
func (s *Session) Back() *Session {
	res, err := s.api.executeRequestVoid(
		http.MethodPost,
		fmt.Sprintf("/session/%s/back", s.id),
	)
	if err != nil {
		handleError(res, err)
	}

	return s
}

// Forward navigates forward in the browser history.
func (s *Session) Forward() *Session {
	res, err := s.api.executeRequestVoid(
		http.MethodPost,
		fmt.Sprintf("/session/%s/forward", s.id),
	)
	if err != nil {
		handleError(res, err)
	}

	return s
}

// GetTitle returns the current page title.
func (s *Session) GetTitle() string {
	res, err := s.api.executeRequestVoid(
		http.MethodGet,
		fmt.Sprintf("/session/%s/title", s.id),
	)
	if err != nil {
		handleError(res, err)
	}

	return res.Value.(string)
}

func (s *Session) GetWindowHandle() string {
	res, err := s.api.executeRequestVoid(
		http.MethodGet,
		fmt.Sprintf("/session/%s/window", s.id),
	)
	if err != nil {
		handleError(res, err)
	}

	return res.Value.(string)
}

func (s *Session) GetWindowHandles() []string {
	res, err := s.api.executeRequestVoid(
		http.MethodGet,
		fmt.Sprintf("/session/%s/window/handles", s.id),
	)
	if err != nil {
		handleError(res, err)
	}

	result := res.Value.([]interface{})
	handles := make([]string, 0, len(result))

	for _, v := range result {
		handles = append(handles, v.(string))
	}

	return handles
}

// If there are no open browsing contexts left, the session is closed.
func (s *Session) CloseWindow() *Session {
	res, err := s.api.executeRequestVoid(
		http.MethodDelete,
		fmt.Sprintf("/session/%s/window", s.id),
	)
	if err != nil {
		handleError(res, err)

		return s
	}

	return s
}

func (s *Session) SwitchHandle(handle string) {
	payload := struct {
		Handle string `json:"handle"`
	}{handle}

	res, err := s.api.executeRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/window", s.id),
		payload,
	)
	if err != nil {
		handleError(res, err)
	}
}

// Opens a new browser tab.
func (s *Session) NewTab() *Handle {
	return s.newWindowWithType(tab)
}

// Opens a new browser tab.
func (s *Session) NewWindow() *Handle {
	return s.newWindowWithType(window)
}

func (s *Session) newWindowWithType(ht handleType) *Handle {
	payload := struct {
		HandleType string `json:"type"`
	}{string(ht)}

	var response struct {
		Value Handle `json:"value"`
	}

	res, err := s.api.executeRequestCustom(
		http.MethodPost,
		fmt.Sprintf("/session/%s/window/new", s.id),
		payload,
		&response,
	)
	if err != nil {
		handleError(
			res,
			errors.Wrapf(err, "failed to open new %s", string(ht)),
		)
	}

	return &response.Value
}

// SwitchToFrame switches current browsing context to the specified iframe using
// the provided element. If nil is provided, the session will switch to the
// top-level browsing context.
func (s *Session) SwitchToFrame(e *element) *Session {
	//nolint:tagliatelle
	type id struct {
		ElementID string `json:"element-6066-11e4-a52e-4f735466cecf"`
	}

	type payload struct {
		ID interface{} `json:"id"`
	}

	var p interface{}

	if e == nil {
		p = payload{nil}
	} else {
		e.setElementID()
		p = payload{id{e.id}}
	}

	res, err := s.api.executeRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/frame", s.id),
		p,
	)
	if err != nil {
		handleError(res, err)
	}

	return s
}

// SwitchToParentFrame switches to the parent frame of the given browsing
// context.
func (s *Session) SwitchToParentFrame() *Session {
	res, err := s.api.executeRequestVoid(
		http.MethodPost,
		fmt.Sprintf("/session/%s/frame/parent", s.id),
	)
	if err != nil {
		handleError(res, err)
	}

	return s
}
