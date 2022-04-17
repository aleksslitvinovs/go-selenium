package selenium

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/logger"
	"github.com/theRealAlpaca/go-selenium/types"
)

// GetText returns the text of the element.
func (e *Element) GetText() string {
	e.setElementID()

	res, err := e.api.executeRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/element/%s/text", e.session.id, e.id),
		e,
	)
	if err != nil {
		handleError(res, err)
	}

	if res.Value == nil {
		return ""
	}

	if v, ok := res.Value.(string); ok {
		return v
	}

	return ""
}

// GetAttribute returns the value of the given attribute of the element. If the
// element does not have the given attribute, an empty string is returned.
func (e *Element) GetAttribute(attribute string) string {
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

	if res.Value == nil {
		logger.Errorf(
			"element %q does not have %q attribute", e.Selector, attribute,
		)

		return ""
	}

	if v, ok := res.Value.(string); ok {
		return v
	}

	return ""
}

// Click clicks on the element.
func (e *Element) Click() *Element {
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

// SendKeys sends the given keys to the element.
func (e *Element) SendKeys(input string) *Element {
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

// Clear clears the text of the element.
func (e *Element) Clear() *Element {
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

// IsPresent checks if the element is present in the DOM.
func (e *Element) IsPresent() bool {
	_, err := e.findElement()
	if err != nil {
		handleError(nil, err)
	}

	return err == nil
}

// IsVisible checks if the element is visible.
func (e *Element) IsVisible() bool {
	return e.handleCondition(
		func() (*response, error) { return e.isVisible() },
	)
}

// IsEnabled checks if the element is enabled.
func (e *Element) IsEnabled() bool {
	return e.handleCondition(
		func() (*response, error) { return e.isEnabled() },
	)
}

// IsSelected checks if the element is selected.
func (e *Element) IsSelected() bool {
	return e.handleCondition(
		func() (*response, error) { return e.isSelected() },
	)
}

func (e *Element) isVisible() (*response, error) {
	e.setElementID()

	return e.api.executeRequest(
		http.MethodGet,
		fmt.Sprintf(
			"/session/%s/element/%s/displayed", e.session.id, e.id,
		),
		e,
	)
}

func (e *Element) isEnabled() (*response, error) {
	e.setElementID()

	return e.api.executeRequest(
		http.MethodGet,
		fmt.Sprintf(
			"/session/%s/element/%s/enabled", e.session.id, e.id,
		),
		e,
	)
}

func (e *Element) isSelected() (*response, error) {
	e.setElementID()

	return e.api.executeRequest(
		http.MethodGet,
		fmt.Sprintf(
			"/session/%s/element/%s/selected", e.session.id, e.id,
		),
		e,
	)
}

func (e *Element) handleCondition(
	condition func() (*response, error),
) bool {
	res, err := condition()
	if err != nil {
		handleError(res, err)

		return false
	}

	if res.Value == nil {
		handleError(nil, errors.New("failed top get element's condition"))

		return false
	}

	if v, ok := res.Value.(bool); ok {
		return v
	}

	return false
}
