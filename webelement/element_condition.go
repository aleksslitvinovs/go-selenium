package webelement

import (
	"fmt"
	"net/http"

	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/util"
)

func (we *webElement) IsPresent() bool {
	_, err := we.findElement()

	return err == nil
}

func (we *webElement) IsVisible() bool {
	return we.handleCondition(
		func() (*api.Response, error) { return we.isVisible() },
	)
}

func (we *webElement) IsEnabled() bool {
	return we.handleCondition(
		func() (*api.Response, error) { return we.isEnabled() },
	)
}

func (we *webElement) IsSelected() bool {
	return we.handleCondition(
		func() (*api.Response, error) { return we.isSelected() },
	)
}

func (we *webElement) isVisible() (*api.Response, error) {
	we.setElementID()

	//nolint:wrapcheck
	return we.api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf(
			"/session/%s/element/%s/displayed", we.session.GetID(), we.id,
		),
		we,
	)
}

func (we *webElement) isEnabled() (*api.Response, error) {
	we.setElementID()

	//nolint:wrapcheck
	return we.api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf(
			"/session/%s/element/%s/enabled", we.session.GetID(), we.id,
		),
		we,
	)
}

func (we *webElement) isSelected() (*api.Response, error) {
	we.setElementID()

	//nolint:wrapcheck
	return we.api.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf(
			"/session/%s/element/%s/selected", we.session.GetID(), we.id,
		),
		we,
	)
}

func (we *webElement) handleCondition(
	condition func() (*api.Response, error),
) bool {
	res, err := condition()
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			util.HandleResponseError(we.session, res.GetErrorReponse())

			return false
		}

		util.HandleError(we.session, err)

		return false
	}

	return res.Value.(bool)
}
