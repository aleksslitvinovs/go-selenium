package webelement

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/util"
)

func (we *webElement) GetText() (string, error) {
	we.setElementID()

	res, err := we.api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/element/%s/text", we.session.GetID(), we.id),
		we,
	)
	if err != nil {
		errRes := res.GetErrorReponse()
		if errRes == nil {
			return "", errors.Wrap(err, "failed to send request to get text")
		}

		util.HandleResponseError(we.session, errRes)
	}

	return res.Value.(string), nil
}

func (we *webElement) Click() error {
	we.setElementID()

	res, err := we.api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element/%s/click", we.session.GetID(), we.id),
		we,
	)
	if err != nil {
		errRes := res.GetErrorReponse()
		if errRes == nil {
			return errors.Wrap(err, "failed to click on element")
		}

		util.HandleResponseError(we.session, errRes)
	}

	return nil
}

func (we *webElement) SendKeys(input string) {
	we.setElementID()

	payload := struct {
		Text string `json:"text"`
	}{input}

	res, err := we.api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element/%s/value", we.session.GetID(), we.id),
		payload,
	)
	if err != nil {
		errRes := res.GetErrorReponse()
		if errRes == nil {
			util.HandleError(
				we.session,
				errors.Wrap(err, "failed to send keys to element"),
			)

			return
		}

		if errRes.Error == api.NoSuchElement && we.settings.IgnoreNotFound {
			return
		}

		util.HandleResponseError(we.session, errRes)
	}
}

func (we *webElement) Clear() error {
	we.setElementID()

	_, err := we.api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element/%s/clear", we.session.GetID(), we.id),
		we,
	)
	if err != nil {
		return errors.Wrap(err, "failed to send request to clear")
	}

	return nil
}

func (we *webElement) setElementID() {
	if we.id != "" {
		return
	}

	intialSettings := *we.settings

	we.settings.IgnoreNotFound = true

	defer func() {
		we.settings = &intialSettings
	}()

	timeout := time.Now().Add(we.settings.RetryTimeout.Duration)

	for time.Now().Before(timeout) {
		if id := we.FindElement(); id != "" {
			we.id = id

			return
		}

		time.Sleep(we.settings.PollInterval.Duration)
	}

	util.HandleError(
		we.session,
		errors.Errorf("Element %q not found", we.Selector),
	)
}
