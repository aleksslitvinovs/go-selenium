package selenium

import (
	"time"

	"github.com/theRealAlpaca/go-selenium/config"
	"github.com/theRealAlpaca/go-selenium/logger"
	"github.com/theRealAlpaca/go-selenium/selector"
	"github.com/theRealAlpaca/go-selenium/types"
	"github.com/theRealAlpaca/go-selenium/webelement"
)

type Element struct {
	SelectorType string                  `json:"using"`
	Selector     string                  `json:"value"`
	Settings     *config.ElementSettings `json:"-"`
	Session      *Session                `json:"-"`
}

var (
	defaultElementSettings = &config.ElementSettings{
		PollInterval: types.Time{Duration: 500 * time.Millisecond},
		RetryTimeout: types.Time{Duration: 5 * time.Second},
		SelectorType: selector.CSS,
	}
)

func (s *Session) NewElement(selector string) types.WebElementer {
	if config.Config.ElementSettings.RetryTimeout.Milliseconds() == 0 {
		logger.Error(`"retry_timeout" must not be 0`)

		s.DeleteSession()
	}

	settings := config.Config.ElementSettings
	if settings == nil {
		settings = defaultElementSettings
	}

	return webelement.NewElement("", s, selector, settings, s.api)
}

func SetSettings(settings *config.ElementSettings) {
	defaultElementSettings = settings
}
