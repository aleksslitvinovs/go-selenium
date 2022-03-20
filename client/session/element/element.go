package element

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/client/session"
)

type Element struct {
	SelectorType string `json:"using"`
	Selector     string `json:"value"`
	webID        string `json:"-"`
}

func NewElement(selectorType, selector string) *Element {
	return &Element{
		SelectorType: selectorType,
		Selector:     selector,
	}
}

func (e *Element) GetText(s *session.Session) (string, error) {
	if err := e.setElementID(s); err != nil {
		return "", errors.Wrap(err, "failed to set element's webID")
	}

	res, err := api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/element/%s/text", s.ID, e.webID),
		s,
		e,
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to send request to get text")
	}

	return res.Value.(string), nil
}

func (e *Element) Click(s *session.Session) error {
	if err := e.setElementID(s); err != nil {
		return errors.Wrap(err, "failed to set element's webID")
	}

	_, err := api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element/%s/click", s.ID, e.webID),
		s,
		e,
	)
	if err != nil {
		return errors.Wrap(err, "failed to send request to click")
	}

	return nil
}

func (e *Element) SendKeys(s *session.Session, input string) error {
	if err := e.setElementID(s); err != nil {
		return errors.Wrap(err, "failed to set element's webID")
	}

	payload := struct {
		Text string `json:"text"`
	}{input}

	_, err := api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element/%s/value", s.ID, e.webID),
		s,
		payload,
	)
	if err != nil {
		return errors.Wrap(err, "failed to send request to send keys")
	}

	return nil
}

func (e *Element) Clear(s *session.Session) error {
	if err := e.setElementID(s); err != nil {
		return errors.Wrap(err, "failed to set element's webID")
	}

	_, err := api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element/%s/clear", s.ID, e.webID),
		s,
		e,
	)
	if err != nil {
		return errors.Wrap(err, "failed to send request to clear")
	}

	return nil
}

func (e *Element) setElementID(s *session.Session) error {
	id, err := e.FindElement(s)
	if err != nil {
		return errors.Wrap(err, "failed to find element")
	}

	e.webID = id

	return nil
}
