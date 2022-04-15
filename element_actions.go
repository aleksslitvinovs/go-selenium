package selenium

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

func (we *Element) GetText() string {
	we.setElementID()

	res, err := we.api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/element/%s/text", we.session.id, we.id),
		we,
	)
	if err != nil {
		errRes := res.GetErrorReponse()
		if errRes == nil {
			HandleError(err)
		}

		HandleResponseError(errRes)
	}

	return res.Value.(string)
}

func (we *Element) Click() *Element {
	we.setElementID()

	res, err := we.api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element/%s/click", we.session.id, we.id),
		we,
	)
	if err != nil {
		errRes := res.GetErrorReponse()
		if errRes == nil {
			HandleError(err)
		}

		HandleResponseError(errRes)
	}

	return we
}

func (we *Element) SendKeys(input string) *Element {
	we.setElementID()

	payload := struct {
		Text string `json:"text"`
	}{input}

	res, err := we.api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element/%s/value", we.session.id, we.id),
		payload,
	)
	if err != nil {
		errRes := res.GetErrorReponse()
		if errRes == nil {
			HandleError(errors.Wrap(err, "failed to send keys to element"))

			return we
		}

		if errRes.Error == NoSuchElement && we.settings.IgnoreNotFound {
			return we
		}

		HandleResponseError(errRes)
	}

	return we
}

func (we *Element) Clear() *Element {
	we.setElementID()

	res, err := we.api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element/%s/clear", we.session.id, we.id),
		we,
	)
	if err != nil {
		errRes := res.GetErrorReponse()
		if errRes == nil {
			HandleError(
				errors.Wrap(err, "failed to send keys to element"),
			)

			return we
		}

		HandleResponseError(errRes)
	}

	return we
}

func (we *Element) IsPresent() bool {
	_, err := we.findElement()

	return err == nil
}

func (we *Element) IsVisible() bool {
	return we.handleCondition(
		func() (*Response, error) { return we.isVisible() },
	)
}

func (we *Element) IsEnabled() bool {
	return we.handleCondition(
		func() (*Response, error) { return we.isEnabled() },
	)
}

func (we *Element) IsSelected() bool {
	return we.handleCondition(
		func() (*Response, error) { return we.isSelected() },
	)
}

func (we *Element) isVisible() (*Response, error) {
	we.setElementID()

	return we.api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf(
			"/session/%s/element/%s/displayed", we.session.id, we.id,
		),
		we,
	)
}

func (we *Element) isEnabled() (*Response, error) {
	we.setElementID()

	return we.api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf(
			"/session/%s/element/%s/enabled", we.session.id, we.id,
		),
		we,
	)
}

func (we *Element) isSelected() (*Response, error) {
	we.setElementID()

	return we.api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf(
			"/session/%s/element/%s/selected", we.session.id, we.id,
		),
		we,
	)
}

func (we *Element) handleCondition(
	condition func() (*Response, error),
) bool {
	res, err := condition()
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			HandleResponseError(res.GetErrorReponse())

			return false
		}

		HandleError(err)

		return false
	}

	return res.Value.(bool)
}
