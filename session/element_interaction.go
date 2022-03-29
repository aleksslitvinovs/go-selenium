package session

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/util"
)

func (e *Element) GetText() (string, error) {
	e.setElementID()

	res, err := api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/element/%s/text", e.Session.ID, e.webID),
		e.Session,
		e,
	)
	if err != nil {
		errRes := res.GetErrorReponse()
		if errRes == nil {
			return "", errors.Wrap(err, "failed to send request to get text")
		}

		util.HandleResponseError(e.Session, errRes)
	}

	return res.Value.(string), nil
}

func (e *Element) Click() error {
	e.setElementID()

	res, err := api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element/%s/click", e.Session.ID, e.webID),
		e.Session,
		e,
	)
	if err != nil {
		errRes := res.GetErrorReponse()
		if errRes == nil {
			return errors.Wrap(err, "failed to click on element")
		}

		util.HandleResponseError(e.Session, errRes)
	}

	return nil
}

func (e *Element) SendKeys(input string) {
	e.setElementID()

	payload := struct {
		Text string `json:"text"`
	}{input}

	res, err := api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element/%s/value", e.Session.ID, e.webID),
		e.Session,
		payload,
	)
	if err != nil {
		errRes := res.GetErrorReponse()
		if errRes == nil {
			util.HandleError(
				e.Session,
				errors.Wrap(err, "failed to send keys to element"),
			)

			return
		}

		if errRes.Error == api.NoSuchElement && e.Settings.IgnoreNotFound {
			return
		}

		util.HandleResponseError(e.Session, errRes)
	}
}

func (e *Element) Clear() error {
	e.setElementID()

	_, err := api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element/%s/clear", e.Session.ID, e.webID),
		e.Session,
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

	intialSettings := *e.Settings

	e.Settings.IgnoreNotFound = true

	defer func() {
		e.Settings = &intialSettings
	}()

	timeout := time.Now().Add(e.Settings.RetryTimeout.Duration)

	for {
		if time.Now().After(timeout) {
			util.HandleError(
				e.Session,
				errors.Errorf("Element %q not found", e.Selector),
			)
		}

		if id := e.FindElement(); id != "" {
			e.webID = id

			return
		}

		time.Sleep(e.Settings.PollInterval.Duration)
	}
}
