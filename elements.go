package selenium

import (
	"fmt"
	"net/http"

	"github.com/aleksslitvinovs/go-selenium/types"
	"github.com/pkg/errors"
)

// Elements describes a collection of Element type.
type Elements struct {
	E

	elements []*Element

	session  *Session
	settings *elementSettings
	api      *apiClient
}

// NewElements return a new Elements. The parameter can be either a selector (
// uses session's default locator) or *E or E struct.
func (s *Session) NewElements(e interface{}) *Elements {
	elem := s.NewElement(e)

	return &Elements{
		E: E{
			Selector:     elem.Selector,
			SelectorType: elem.SelectorType,
		},
		session:  s,
		settings: config.Element,
		api:      s.api,
	}
}

// Size returns the number of elements.
func (ee *Elements) Size() int {
	ee.setElementsID()

	return len(ee.elements)
}

// Elements returns the list elements.
func (ee *Elements) Elements() []*Element {
	return ee.elements
}

func (ee *Elements) findElements() ([]string, error) {
	res, err := ee.api.executeRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/elements", ee.session.id),
		ee.E,
	)
	if err != nil {
		errRes := res.getErrorReponse()
		if errRes == nil {
			return []string{}, errors.Wrap(err, "failed to find element")
		}

		if errors.As(errRes, &types.ErrNoSuchElement) &&
			ee.settings.IgnoreNotFound {
			return []string{}, nil
		}

		ok := isAllowedError(errRes)
		if !ok {
			panic(errRes.String())
		}
	}

	v, ok := res.Value.([]interface{})
	if !ok {
		return []string{}, errors.New("failed to convert GET elements response")
	}

	ids := make([]string, 0, len(v))

	for _, elem := range v {
		e, ok := elem.(map[string]interface{})

		if !ok {
			return []string{}, errors.New(
				"failed to convert elements' element ID response",
			)
		}

		ids = append(ids, getElementID(e))
	}

	return ids, nil
}

func (ee *Elements) setElementsID() {
	if len(ee.elements) != 0 {
		return
	}

	ids, err := ee.findElements()
	if err != nil {
		handleError(nil, err)
	}

	for _, id := range ids {
		e := ee.session.NewElement(id)
		e.id = id

		ee.elements = append(ee.elements, e)
	}
}
