package session

import (
	"fmt"
	"net/http"

	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/util"
)

func (e *Element) IsVisible() bool {
	return e.handleCondition(
		func() (*api.Response, error) { return e.isVisible() },
	)
}

func (e *Element) IsEnabled() bool {
	return e.handleCondition(
		func() (*api.Response, error) { return e.isEnabled() },
	)
}

func (e *Element) IsSelected() bool {
	return e.handleCondition(
		func() (*api.Response, error) { return e.isSelected() },
	)
}

func (e *Element) isVisible() (*api.Response, error) {
	e.setElementID()

	if e.webID == "" {
		return &api.Response{}, ErrWebIDNotSet
	}

	//nolint:wrapcheck
	return api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/element/%s/displayed", e.Session.ID, e.webID),
		e.Session,
		e,
	)
}

func (e *Element) isEnabled() (*api.Response, error) {
	e.setElementID()

	if e.webID == "" {
		return &api.Response{}, ErrWebIDNotSet
	}

	//nolint:wrapcheck
	return api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/element/%s/enabled", e.Session.ID, e.webID),
		e.Session,
		e,
	)
}

func (e *Element) isSelected() (*api.Response, error) {
	e.setElementID()

	if e.webID == "" {
		return &api.Response{}, ErrWebIDNotSet
	}

	//nolint:wrapcheck
	return api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/element/%s/selected", e.Session.ID, e.webID),
		e.Session,
		e,
	)
}

func (e *Element) handleCondition(
	condition func() (*api.Response, error),
) bool {
	res, err := condition()
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			util.HandleResponseError(e.Session, res.GetErrorReponse())

			return false
		}

		util.HandleError(e.Session, err)

		return false
	}

	return res.Value.(bool)
}
