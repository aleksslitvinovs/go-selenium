package selenium

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/types"
)

// Elements describes a collection of Element type.
type Elements struct {
	E

	elements []string `json:"-"`

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
		return []string{}, errors.New("failed to convert elements' ID response")
	}

	fmt.Printf("Value type %T\n", v)
	fmt.Printf("Value value %v\n", len(v))

	ids := make([]string, 0, len(v))

	for _, elem := range v {
		e, ok := elem.(map[string]string)
		if !ok {
			return []string{}, errors.New(
				"failed to convert elements' element ID response",
			)
		}

		// TODO: Make this actually work
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

	ee.elements = ids
}
