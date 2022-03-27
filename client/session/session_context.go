package session

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/util"
)

type Contexter interface {
	// GetWindowHandle returns the current browsing context window handle.
	GetWindowHandle() string

	// GetWindowHandles returns an array of window handles for each browsing
	// context. The order of the handles is arbitrary.
	GetWindowHandles() []string

	// CloseWindow closes the current browsing context.
	CloseWindow()

	// SwitchHandle switches to the given browsing context.
	SwitchHandle(handle string)

	// NewTab opens a new tab in the current window.
	NewTab() *Handle

	// NewWindow opens a window.
	NewWindow() *Handle

	// TODO: Not working yet.
	// SwitchToFrame switches to the given iframe element.
	SwitchToFrame(id int)

	// SwitchToParentFrame switches to the parent frame of the current frame.
	SwitchToParentFrame()
}

type Handle struct {
	ID   string `json:"handle"`
	Type string `json:"type"`
}

type handleType string

var (
	tab    handleType = "tab"
	window handleType = "window"
)

func (s *Session) GetWindowHandle() string {
	res, err := api.ExecuteRequestVoid(
		http.MethodGet,
		fmt.Sprintf("/session/%s/window", s.ID),
		s,
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			util.HandleResponseError(s, errRes)

			return ""
		}

		util.HandleError(s, errors.Wrap(err, "failed to get window handle"))
	}

	return res.Value.(string)
}

func (s *Session) GetWindowHandles() []string {
	res, err := api.ExecuteRequestVoid(
		http.MethodGet,
		fmt.Sprintf("/session/%s/window/handles", s.ID),
		s,
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			util.HandleResponseError(s, errRes)

			return []string{}
		}

		util.HandleError(s, errors.Wrap(err, "failed to get window handles"))
	}

	result := res.Value.([]interface{})
	handles := make([]string, 0, len(result))

	for _, v := range result {
		handles = append(handles, v.(string))
	}

	return handles
}

func (s *Session) CloseWindow() {
	res, err := api.ExecuteRequestVoid(
		http.MethodDelete,
		fmt.Sprintf("/session/%s/window", s.ID),
		s,
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			util.HandleResponseError(s, errRes)

			return
		}

		util.HandleError(s, errors.Wrap(err, "failed to close window"))
	}

	s.KillDriver <- struct{}{}
}

func (s *Session) SwitchHandle(handle string) {
	payload := struct {
		Handle string `json:"handle"`
	}{handle}

	res, err := api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/window", s.ID),
		s,
		payload,
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			util.HandleResponseError(s, errRes)

			return
		}

		util.HandleError(s, errors.Wrap(err, "failed to switch to handle"))
	}
}

func (s *Session) NewTab() *Handle {
	return s.newWindowWithType(tab)
}

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

	err := api.ExecuteRequestCustom(
		http.MethodPost,
		fmt.Sprintf("/session/%s/window/new", s.ID),
		s,
		payload,
		&response,
	)
	if err != nil {
		util.HandleError(
			s, errors.Wrapf(err, "failed to open new %s", string(ht)),
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

	res, err := api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/frame", s.ID),
		s,
		payload,
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			util.HandleResponseError(s, errRes)

			return
		}

		util.HandleError(s, errors.Wrap(err, "failed to get window handles"))
	}
}

func (s *Session) SwitchToParentFrame() {
	res, err := api.ExecuteRequestVoid(
		http.MethodPost,
		fmt.Sprintf("/session/%s/frame/parent", s.ID),
		s,
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			util.HandleResponseError(s, errRes)

			return
		}

		util.HandleError(s, errors.Wrap(err, "failed to get window handles"))
	}
}
