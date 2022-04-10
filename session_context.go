package selenium

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/types"
	"github.com/theRealAlpaca/go-selenium/util"
)

type handleType string

var (
	tab    handleType = "tab"
	window handleType = "window"
)

func (s *Session) GetWindowHandle() string {
	res, err := s.api.ExecuteRequestVoid(
		http.MethodGet,
		fmt.Sprintf("/session/%s/window", s.id),
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
	res, err := s.api.ExecuteRequestVoid(
		http.MethodGet,
		fmt.Sprintf("/session/%s/window/handles", s.id),
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
	res, err := s.api.ExecuteRequestVoid(
		http.MethodDelete,
		fmt.Sprintf("/session/%s/window", s.id),
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			util.HandleResponseError(s, errRes)

			return
		}

		util.HandleError(s, errors.Wrap(err, "failed to close window"))
	}

	handles := s.GetWindowHandles()

	// If there are no open browsing contexts left, the session is closed.
	if len(handles) == 0 {
		s.killDriver <- struct{}{}
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
			util.HandleResponseError(s, errRes)

			return
		}

		util.HandleError(s, errors.Wrap(err, "failed to switch to handle"))
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

	res, err := s.api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/frame", s.id),
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
	res, err := s.api.ExecuteRequestVoid(
		http.MethodPost,
		fmt.Sprintf("/session/%s/frame/parent", s.id),
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			util.HandleResponseError(s, errRes)

			return
		}

		util.HandleError(s, errors.Wrap(err, "failed to get window handles"))
	}
}
