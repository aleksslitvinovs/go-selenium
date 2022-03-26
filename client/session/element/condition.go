package element

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/client/session"
	"github.com/theRealAlpaca/go-selenium/util"
)

func (e *Element) IsVisible(s *session.Session) (bool, error) {
	if err := e.setElementID(); err != nil {
		return false, errors.Wrap(err, "failed to set element's webID")
	}

	res, err := api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/element/%s/displayed", s.ID, e.webID),
		s,
		e,
	)
	if err != nil {
		util.HandleResponseError(s, res.GetErrorReponse())

		return false, nil
	}

	return res.Value.(bool), nil
}

func (e *Element) IsEnabled(s *session.Session) (bool, error) {
	if err := e.setElementID(); err != nil {
		return false, errors.Wrap(err, "failed to set element's webID")
	}

	res, err := api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/element/%s/enabled", s.ID, e.webID),
		s,
		e,
	)
	if err != nil {
		util.HandleResponseError(s, res.GetErrorReponse())

		return false, nil
	}

	return res.Value.(bool), nil
}

func (e *Element) IsSelected(s *session.Session) (bool, error) {
	if err := e.setElementID(); err != nil {
		return false, errors.Wrap(err, "failed to set element's webID")
	}

	res, err := api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/element/%s/selected", s.ID, e.webID),
		s,
		e,
	)
	if err != nil {
		util.HandleResponseError(s, res.GetErrorReponse())

		return false, nil
	}

	return res.Value.(bool), nil
}
