package config

import (
	"encoding/json"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/logger"
	"github.com/theRealAlpaca/go-selenium/selector"
	"github.com/theRealAlpaca/go-selenium/types"
)

//nolint:tagliatelle
type ElementSettings struct {
	IgnoreNotFound bool       `json:"ignore_not_found"`
	RetryTimeout   types.Time `json:"retry_timeout"`
	PollInterval   types.Time `json:"poll_interval"`
	SelectorType   string     `json:"selector_type"`
}

//nolint:tagliatelle
type WebDriverConfig struct {
	AutoStart    bool                   `json:"auto_start"`
	PathToBinary string                 `json:"path"`
	URL          string                 `json:"url"`
	Timeout      types.Time             `json:"timeout"`
	Capabalities map[string]interface{} `json:"capabilities"`
}

//nolint:tagliatelle
type SessionSettings struct {
	AutoStart bool `json:"auto_start"`
}

//nolint:tagliatelle
type config struct {
	LogLevel                 logger.LevelName `json:"logging"`
	SoftAsserts              bool             `json:"soft_asserts"`
	RaiseErrorsAutomatically bool             `json:"raise_errors_automatically"` //nolint:lll
	ElementSettings          *ElementSettings `json:"element_settings,omitempty"` //nolint:lll
	Session                  *SessionSettings `json:"session,omitempty"`
	// TODO: Allow running multiple drivers.
	WebDriver *WebDriverConfig `json:"webdriver,omitempty"`
}

var Config = &config{LogLevel: "info"}

const defaultConfigPath = "goseleniumrc.json"

// ReadConfig reads the config file and returns a Config struct.
func ReadConfig(configPath string) error {
	if configPath == "" {
		configPath = defaultConfigPath
	}

	_, err := os.Stat(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logger.Info(
				"No config file found. Will create and use default config.",
			)

			c, err := createDefaultConfig()
			if err != nil {
				panic(errors.Wrap(err, "failed to create default config"))
			}

			Config = c

			return nil
		}

		panic(err)
	}

	c, err := readConfigFromFile(configPath)
	if err != nil {
		panic(err)
	}

	c.validateConfig()

	if err := c.writeToConfig(configPath); err != nil {
		panic(err)
	}

	Config = c

	return nil
}

func readConfigFromFile(configPath string) (*config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config, errors.Wrap(err, "failed to read config file")
	}

	var c config

	if err := json.Unmarshal(data, &c); err != nil {
		return Config, errors.Wrap(err, "failed to parse config file")
	}

	return &c, nil
}

func createDefaultConfig() (*config, error) {
	// TODO: implement automatic driver download.

	c := &config{
		LogLevel:                 logger.InfoLvl,
		SoftAsserts:              false,
		RaiseErrorsAutomatically: true,
		// ElementSettings:          &ElementSettings{},
		// WebDriver:                &WebDriverConfig{},
	}

	err := c.writeToConfig(defaultConfigPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to write config")
	}

	return c, nil
}

func (c *config) validateConfig() {
	if c.ElementSettings.SelectorType == "" {
		logger.Warn(`"selector_type" is not set. Defaulting to "css".`)

		c.ElementSettings.SelectorType = selector.CSS
	}

	if c.ElementSettings.RetryTimeout.Duration == 0 {
		logger.Warn(`"retry_timeout" is not set. Defaulting to "10s".`)

		c.ElementSettings.RetryTimeout = types.Time{Duration: 10 * time.Second}
	}

	if c.ElementSettings.PollInterval.Duration == 0 {
		logger.Warn(`"poll_interval" is not set. Defaulting to "500ms".`)

		c.ElementSettings.PollInterval = types.Time{
			Duration: 500 * time.Millisecond,
		}
	}

	if c.WebDriver.Timeout.Duration == 0 {
		logger.Warn(`"timeout" is not set. Defaulting to "10s".`)

		c.WebDriver.Timeout = types.Time{Duration: 10 * time.Second}
	}

	if c.WebDriver.URL == "" {
		logger.Warn(`"url" is not set. Defaulting to "http://localhost:4444".`)

		c.WebDriver.URL = "http://localhost:4444"
	}
}

func (c *config) writeToConfig(configPath string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal config")
	}

	f, err := os.Create(configPath)
	if err != nil {
		return errors.Wrap(err, "failed to create config file")
	}

	defer f.Close()

	_, err = f.WriteString(string(data))
	if err != nil {
		return errors.Wrap(err, "failed to write config")
	}

	return nil
}
