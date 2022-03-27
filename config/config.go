package config

import (
	"encoding/json"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/logger"
)

//nolint:tagliatelle
type ElementSettings struct {
	IgnoreNotFound bool          `json:"ignore_not_found"`
	RetryTimeout   time.Duration `json:"retry_timeout"`
	PollInterval   time.Duration `json:"poll_interval"`
	SelectorType   string        `json:"selector_type"`
}

type WebDriverConfig struct {
	PathToBinary string        `json:"path"`
	URL          string        `json:"url"`
	Port         int           `json:"port"`
	Timeout      time.Duration `json:"timeout"`
}

//nolint:tagliatelle
type config struct {
	LogLevel                 logger.LevelName `json:"logging"`
	SoftAsserts              bool             `json:"soft_asserts"`
	RaiseErrorsAutomatically bool             `json:"raise_errors_automatically"` //nolint:lll
	ElementSettings          *ElementSettings `json:"element_settings,omitempty"` //nolint:lll

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
	f, err := os.Create(defaultConfigPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create config file")
	}

	defer f.Close()

	// TODO: implement automatic driver download.

	c := &config{
		LogLevel:                 logger.InfoLvl,
		SoftAsserts:              false,
		RaiseErrorsAutomatically: true,
		// ElementSettings:          &ElementSettings{},
		// WebDriver:                &WebDriverConfig{},
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal config")
	}

	_, err = f.WriteString(string(data))
	if err != nil {
		return nil, errors.Wrap(err, "failed to write config")
	}

	return c, nil
}
