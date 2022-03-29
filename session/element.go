package session

import (
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/config"
	"github.com/theRealAlpaca/go-selenium/logger"
	"github.com/theRealAlpaca/go-selenium/selector"
)

type Element struct {
	SelectorType string                  `json:"using"`
	Selector     string                  `json:"value"`
	Settings     *config.ElementSettings `json:"-"`
	Session      *Session                `json:"-"`
	webID        string                  `json:"-"`
}

var (
	ErrWebIDNotSet = errors.New("WebID not set")

	defaultElementSettings = &config.ElementSettings{
		PollInterval: config.Time{Duration: 500 * time.Millisecond},
		RetryTimeout: config.Time{Duration: 5 * time.Second},
		SelectorType: selector.CSS,
	}
)

func (s *Session) NewElement(selector string) *Element {
	if config.Config.ElementSettings.RetryTimeout.Milliseconds() == 0 {
		logger.Error(`"retry_timeout" must not be 0`)

		s.DeleteSession()
	}

	settings := config.Config.ElementSettings
	if settings == nil {
		settings = defaultElementSettings
	}

	return &Element{
		SelectorType: settings.SelectorType,
		Selector:     selector,
		Settings:     settings,
		Session:      s,
	}
}

func SetSettings(settings *config.ElementSettings) {
	defaultElementSettings = settings
}

func UseCSS() {
	defaultElementSettings.SelectorType = selector.CSS
}

func UseXPath() {
	defaultElementSettings.SelectorType = selector.XPath
}
