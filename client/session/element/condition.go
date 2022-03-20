package element

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/client/session"
)

func (e *Element) IsVisible(s *session.Session) (bool, error) {
	if err := e.setElementID(s); err != nil {
		return false, errors.Wrap(err, "failed to set element's webID")
	}

	res, err := api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/element/%s/displayed", s.ID, e.webID),
		s,
		e,
	)
	if err != nil {
		return false, errors.Wrap(err, "could not get display state")
	}

	return res.Value.(bool), nil
}

func (e *Element) IsEnabled(s *session.Session) (bool, error) {
	if err := e.setElementID(s); err != nil {
		return false, errors.Wrap(err, "failed to set element's webID")
	}

	res, err := api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/element/%s/enabled", s.ID, e.webID),
		s,
		e,
	)
	if err != nil {
		return false, errors.Wrap(err, "could not get enabled stated")
	}

	return res.Value.(bool), nil
}

func (e *Element) IsSelected(s *session.Session) (bool, error) {
	if err := e.setElementID(s); err != nil {
		return false, errors.Wrap(err, "failed to set element's webID")
	}

	res, err := api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/element/%s/selected", s.ID, e.webID),
		s,
		e,
	)
	if err != nil {
		return false, errors.Wrap(err, "could not get selected state")
	}

	return res.Value.(bool), nil
}
