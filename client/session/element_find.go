package session

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/util"
)

func (e *Element) FindElement() string {
	res, err := api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element", e.Session.ID),
		e.Session,
		e,
	)
	if err != nil {
		errRes := res.GetErrorReponse()

		if errRes != nil {
			if errRes.Error == api.NoSuchElement && e.Settings.IgnoreNotFound {
				return ""
			}

			util.HandleResponseError(e.Session, errRes)

			return ""
		}

		util.HandleError(e.Session, err)

		return ""
	}

	v, ok := res.Value.(map[string]string)
	if !ok {
		util.HandleError(e.Session, errors.New("failed to find element"))

		return ""
	}

	for _, v := range v {
		if v != "" {
			return v
		}
	}

	util.HandleError(e.Session, errors.New("failed to get element id"))

	return ""
}
