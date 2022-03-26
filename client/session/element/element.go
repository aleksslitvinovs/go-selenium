package element

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/client/session"
	"github.com/theRealAlpaca/go-selenium/util"
)

type Element struct {
	SelectorType   string           `json:"using"`
	Selector       string           `json:"value"`
	IgnoreNotFound bool             `json:"-"`
	Session        *session.Session `json:"-"`
	webID          string           `json:"-"`
}

func NewElement(s *session.Session, selectorType, selector string) *Element {
	return &Element{
		Session:      s,
		SelectorType: selectorType,
		Selector:     selector,
	}
}

func (e *Element) GetText(s *session.Session) (string, error) {
	e.setElementID()

	res, err := api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/element/%s/text", s.ID, e.webID),
		s,
		e,
	)
	if err != nil {
		errRes := res.GetErrorReponse()
		if errRes == nil {
			return "", errors.Wrap(err, "failed to send request to get text")
		}

		util.HandleResponseError(s, errRes)
	}

	return res.Value.(string), nil
}

func (e *Element) Click(s *session.Session) error {
	e.setElementID()

	res, err := api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element/%s/click", s.ID, e.webID),
		s,
		e,
	)
	if err != nil {
		errRes := res.GetErrorReponse()
		if errRes == nil {
			return errors.Wrap(err, "failed to click on element")
		}

		util.HandleResponseError(s, errRes)
	}

	return nil
}

func (e *Element) SendKeys(s *session.Session, input string) {
	e.setElementID()

	payload := struct {
		Text string `json:"text"`
	}{input}

	res, err := api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element/%s/value", s.ID, e.webID),
		s,
		payload,
	)
	if err != nil {
		errRes := res.GetErrorReponse()
		if errRes == nil {
			util.HandleError(
				s, errors.Wrap(err, "failed to send keys to element"),
			)

			return
		}

		if errRes.Error == api.NoSuchElement && e.IgnoreNotFound {
			return
		}

		util.HandleResponseError(s, errRes)
	}
}

func (e *Element) Clear(s *session.Session) error {
	e.setElementID()

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

func (e *Element) setElementID() {
	if e.webID != "" {
		return
	}

	id, err := e.FindElement()
	if err != nil {
		util.HandleError(e.Session, err)
	}

	e.webID = id
}
