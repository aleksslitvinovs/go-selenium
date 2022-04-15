package webelement

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/types"
	"github.com/theRealAlpaca/go-selenium/util"
)

func (we *webElement) GetText() string {
	we.setElementID()

	res, err := we.api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/element/%s/text", we.session.GetID(), we.id),
		we,
	)
	if err != nil {
		errRes := res.GetErrorReponse()
		if errRes == nil {
			util.HandleError(we.session, err)
		}

		util.HandleResponseError(we.session, errRes)
	}

	return res.Value.(string)
}

func (we *webElement) Click() types.WebElementer {
	we.setElementID()

	res, err := we.api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element/%s/click", we.session.GetID(), we.id),
		we,
	)
	if err != nil {
		errRes := res.GetErrorReponse()
		if errRes == nil {
			util.HandleError(we.session, err)
		}

		util.HandleResponseError(we.session, errRes)
	}

	return we
}

func (we *webElement) SendKeys(input string) types.WebElementer {
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

			return we
		}

		if errRes.Error == api.NoSuchElement && we.settings.IgnoreNotFound {
			return we
		}

		util.HandleResponseError(we.session, errRes)
	}

	return we
}

func (we *webElement) Clear() types.WebElementer {
	we.setElementID()

	res, err := we.api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element/%s/clear", we.session.GetID(), we.id),
		we,
	)
	if err != nil {
		errRes := res.GetErrorReponse()
		if errRes == nil {
			util.HandleError(
				we.session,
				errors.Wrap(err, "failed to send keys to element"),
			)

			return we
		}

		util.HandleResponseError(we.session, errRes)
	}

	return we
}
