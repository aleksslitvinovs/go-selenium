package webelement

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/util"
)

const (
	// https://www.w3.org/TR/webdriver/#elements
	webElementID    = "element-6066-11e4-a52e-4f735466cecf"
	legacyElementID = "ELEMENT"
)

func (we *webElement) FindElement() string {
	res, err := we.api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element", we.session.GetID()),
		we,
	)
	if err != nil {
		errRes := res.GetErrorReponse()

		if errRes != nil {
			if errRes.Error == api.NoSuchElement && we.settings.IgnoreNotFound {
				return ""
			}

			util.HandleResponseError(we.session, errRes)

			return ""
		}

		util.HandleError(we.session, err)

		return ""
	}

	v, ok := res.Value.(map[string]string)
	if !ok {
		util.HandleError(we.session, errors.New("failed to find element"))

		return ""
	}

	id := getElementID(v)

	if id == "" {
		util.HandleError(we.session, errors.New("failed to get element id"))

		return ""
	}

	return id
}

func getElementID(elements map[string]string) string {
	supportedIDs := []string{webElementID, legacyElementID}

	for _, key := range supportedIDs {
		e, ok := elements[key]
		if !ok || e == "" {
			continue
		}

		return e
	}

	return ""
}
