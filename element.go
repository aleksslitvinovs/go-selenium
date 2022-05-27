package selenium

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/logger"
	"github.com/theRealAlpaca/go-selenium/types"
)

// E is a helper struct that represents web element.
type E struct {
	Selector     string `json:"value"`
	SelectorType string `json:"using"`
}

// Element describes web element.
type Element struct {
	E

	id       string
	session  *Session
	settings *elementSettings
	api      *apiClient
}

const (
	// Based on https://www.w3.org/TR/webdriver/#elements
	webElementID    = "element-6066-11e4-a52e-4f735466cecf"
	legacyElementID = "ELEMENT"
)

// NewElement returns a new Element. The parameter can be either a selector (
// uses session's default locator) or *E or E struct.
func (s *Session) NewElement(e interface{}) *Element {
	if e == nil {
		return nil
	}

	switch v := e.(type) {
	case string:
		return &Element{
			E: E{
				Selector:     v,
				SelectorType: s.locatorStrategy,
			},
			settings: config.Element,
			session:  s,
			api:      s.api,
		}
	case *E:
		if v.SelectorType == "" {
			v.SelectorType = s.locatorStrategy
		}

		return &Element{
			E:        *v,
			settings: config.Element,
			session:  s,
			api:      s.api,
		}
	case E:
		if v.SelectorType == "" {
			v.SelectorType = s.locatorStrategy
		}

		return &Element{
			E:        v,
			settings: config.Element,
			session:  s,
			api:      s.api,
		}
	default:
		panic(errors.Errorf("unsupported element type: %T", v))
	}
}

func (e *Element) setElementID() {
	if e.id != "" {
		return
	}

	intialSettings := *e.settings

	e.settings.IgnoreNotFound = true

	defer func() {
		e.settings = &intialSettings
	}()

	timeout := time.Now().Add(e.settings.RetryTimeout.Duration)

	var err error

	for time.Now().Before(timeout) {
		id, err := e.findElement()
		if err != nil {
			handleError(nil, err)
		}

		if id == "" {
			time.Sleep(e.settings.PollInterval.Duration)

			continue
		}

		e.id = id

		return
	}

	if err != nil {
		logger.Debugf("An error occurred while finding element: %s", err)
	}

	handleError(
		nil,
		errors.Errorf(
			"Element %q (%s) not found", e.Selector, e.SelectorType,
		),
	)
}

func (e *Element) findElement() (string, error) {
	res, err := e.api.executeRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element", e.session.id),
		e,
	)
	if err != nil {
		errRes := res.getErrorReponse()
		if errRes == nil {
			return "", errors.Wrap(err, "failed to find element")
		}

		if errors.As(errRes, &types.ErrNoSuchElement) &&
			e.settings.IgnoreNotFound {
			return "", nil
		}

		ok := isAllowedError(errRes)
		if !ok {
			panic(errRes.String())
		}
	}

	v, ok := res.Value.(map[string]string)
	if !ok {
		return "", errors.New("failed to convert element's ID response")
	}

	id := getElementID(v)

	if id == "" {
		return "", errors.New("failed to get element id")
	}

	return id, nil
}

func isAllowedError(err error) bool {
	if errors.As(err, &types.ErrStaleElementReference) {
		return true
	}

	if errors.As(err, &types.ErrElementlickIntercepted) {
		return true
	}

	if errors.As(err, &types.ErrElementNotInteractable) {
		return true
	}

	return false
}

func getElementID[T string | interface{}](elements map[string]T) string {
	supportedIDs := []string{webElementID, legacyElementID}

	for _, key := range supportedIDs {
		e, ok := elements[key]
		if !ok {
			continue
		}

		id, ok := any(e).(string)
		if !ok {
			continue
		}

		return id
	}

	return ""
}
