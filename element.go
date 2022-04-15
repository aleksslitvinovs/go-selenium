package selenium

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/logger"
	"github.com/theRealAlpaca/go-selenium/selector"
	"github.com/theRealAlpaca/go-selenium/types"
)

type Element struct {
	Selector     string `json:"value"`
	SelectorType string `json:"using"`

	id       string
	session  *Session
	settings *ElementSettings
	api      *APIClient
}

var defaultSettings = &ElementSettings{
	PollInterval: types.Time{Duration: 500 * time.Millisecond},
	RetryTimeout: types.Time{Duration: 5 * time.Second},
	SelectorType: selector.CSS,
}

const (
	// Based on https://www.w3.org/TR/webdriver/#elements
	webElementID    = "element-6066-11e4-a52e-4f735466cecf"
	legacyElementID = "ELEMENT"
)

func (s *Session) NewElement(selector string) *Element {
	return &Element{
		Selector:     selector,
		SelectorType: s.defaultLocator,
		settings:     defaultSettings,
		session:      s,
		api:          s.api,
	}
}

func (we *Element) setElementID() {
	if we.id != "" {
		return
	}

	intialSettings := *we.settings

	we.settings.IgnoreNotFound = true

	defer func() {
		we.settings = &intialSettings
	}()

	timeout := time.Now().Add(we.settings.RetryTimeout.Duration)

	var err error

	for time.Now().Before(timeout) {
		id, err := we.findElement()
		if err != nil || id == "" {
			time.Sleep(we.settings.PollInterval.Duration)

			continue
		}

		we.id = id

		return
	}

	if err != nil {
		logger.Debugf("An error occurred while finding element: %s", err)
	}

	HandleError(
		errors.Errorf(
			"Element %q (%s) not found", we.Selector, we.SelectorType,
		),
	)
}

func (we *Element) findElement() (string, error) {
	res, err := we.api.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element", we.session.GetID()),
		we,
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to find element")
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
