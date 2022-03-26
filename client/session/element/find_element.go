package element

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/util"
)

func (e *Element) FindElement() (string, error) {
	res, err := api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element", e.Session.ID),
		e.Session,
		e,
	)
	if err != nil {
		errRes := res.GetErrorReponse()

		if errRes == nil {
			return "", errors.Wrap(err, "failed to find element")
		}

		if errRes.Error == api.NoSuchElement && e.IgnoreNotFound {
			return "", nil
		}

		util.HandleResponseError(e.Session, errRes)
	}

	v, ok := res.Value.(map[string]string)
	if !ok {
		return "", errors.New("failed to find element")
	}

	for _, v := range v {
		if v != "" {
			return v, nil
		}
	}

	return "", errors.New("failed to get element id")
}
