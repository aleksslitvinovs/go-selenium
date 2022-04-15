package config

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/logger"
	"github.com/theRealAlpaca/go-selenium/types"
)

//nolint:tagliatelle
type RunnerSettings struct {
	ParallelRuns int `json:"parallel_runs"`
}

//nolint:tagliatelle
type ElementSettings struct {
	IgnoreNotFound bool       `json:"ignore_not_found"`
	RetryTimeout   types.Time `json:"retry_timeout"`
	PollInterval   types.Time `json:"poll_interval"`
	SelectorType   string     `json:"selector_type"`
}

//nolint:tagliatelle
type WebDriverConfig struct {
	ManualStart  bool                   `json:"manual_start"`
	PathToBinary string                 `json:"path"`
	URL          string                 `json:"url"`
	Timeout      types.Time             `json:"timeout"`
	Capabalities map[string]interface{} `json:"capabilities"`
}

//nolint:tagliatelle
type config struct {
	LogLevel                 logger.LevelName `json:"logging"`
	SoftAsserts              bool             `json:"soft_asserts"`
	Runner                   *RunnerSettings  `json:"runner"`
	RaiseErrorsAutomatically bool             `json:"raise_errors_automatically"` //nolint:lll
	ElementSettings          *ElementSettings `json:"element_settings,omitempty"` //nolint:lll
	// TODO: Allow running multiple drivers.
	WebDriver *WebDriverConfig `json:"webdriver,omitempty"`
}

var Config *config

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
				return errors.Wrap(err, "failed to create default config")
			}

			Config = c

			return nil
		}

		return errors.Wrap(err, "failed to stat config file")
	}

	c, err := readConfigFromFile(configPath)
	if err != nil {
		return errors.Wrap(err, "failed to read config file")
	}

	c.validateConfig()

	if err := c.writeToConfig(configPath); err != nil {
		return errors.Wrap(err, "failed to write config")
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
	c := &config{
		LogLevel:                 logger.InfoLvl,
		SoftAsserts:              false,
		Runner:                   &RunnerSettings{ParallelRuns: 1},
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
