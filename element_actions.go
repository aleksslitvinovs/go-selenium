package selenium

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/types"
)

func (e *element) GetText() string {
	e.setElementID()

	res, err := e.api.executeRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/element/%s/text", e.session.id, e.id),
		e,
	)
	if err != nil {
		handleError(res, err)
	}

	return res.Value.(string)
}

func (e *element) GetAttribute(attribute string) string {
	e.setElementID()

	res, err := e.api.executeRequestVoid(
		http.MethodGet,
		fmt.Sprintf(
			"/session/%s/element/%s/attribute/%s",
			e.session.id, e.id, attribute,
		),
	)
	if err != nil {
		handleError(res, err)
	}

	return res.Value.(string)
}

func (e *element) Click() *element {
	e.setElementID()

	res, err := e.api.executeRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element/%s/click", e.session.id, e.id),
		e,
	)
	if err != nil {
		handleError(res, err)
	}

	return e
}

func (e *element) SendKeys(input string) *element {
	e.setElementID()

	payload := struct {
		Text string `json:"text"`
	}{input}

	res, err := e.api.executeRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element/%s/value", e.session.id, e.id),
		payload,
	)
	if err != nil {
		errRes := res.getErrorReponse()
		if errRes == nil {
			handleError(nil, err)

			return e
		}

		if errors.As(errRes, &types.ErrNoSuchElement) &&
			e.settings.IgnoreNotFound {
			return e
		}

		handleError(res, err)
	}

	return e
}

func (e *element) Clear() *element {
	e.setElementID()

	res, err := e.api.executeRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element/%s/clear", e.session.id, e.id),
		e,
	)
	if err != nil {
		handleError(res, err)
	}

	return e
}

func (e *element) IsPresent() bool {
	_, err := e.findElement()
	if err != nil {
		handleError(nil, err)
	}

	return err == nil
}

func (e *element) IsVisible() bool {
	return e.handleCondition(
		func() (*response, error) { return e.isVisible() },
	)
}

func (e *element) IsEnabled() bool {
	return e.handleCondition(
		func() (*response, error) { return e.isEnabled() },
	)
}

func (e *element) IsSelected() bool {
	return e.handleCondition(
		func() (*response, error) { return e.isSelected() },
	)
}

func (e *element) isVisible() (*response, error) {
	e.setElementID()

	return e.api.executeRequest(
		http.MethodGet,
		fmt.Sprintf(
			"/session/%s/element/%s/displayed", e.session.id, e.id,
		),
		e,
	)
}

func (e *element) isEnabled() (*response, error) {
	e.setElementID()

	return e.api.executeRequest(
		http.MethodGet,
		fmt.Sprintf(
			"/session/%s/element/%s/enabled", e.session.id, e.id,
		),
		e,
	)
}

func (e *element) isSelected() (*response, error) {
	e.setElementID()

	return e.api.executeRequest(
		http.MethodGet,
		fmt.Sprintf(
			"/session/%s/element/%s/selected", e.session.id, e.id,
		),
		e,
	)
}

func (e *element) handleCondition(
	condition func() (*response, error),
) bool {
	res, err := condition()
	if err != nil {
		handleError(res, err)

		return false
	}

	return res.Value.(bool)
}
