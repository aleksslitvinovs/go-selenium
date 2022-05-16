package selenium

import (
	"encoding/json"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/logger"
	"github.com/theRealAlpaca/go-selenium/selector"
	"github.com/theRealAlpaca/go-selenium/types"
)

type runnerSettings struct {
	ParallelRuns int `json:"parallel_runs"`
}

type elementSettings struct {
	SelectorType   string     `json:"selector_type"`
	IgnoreNotFound bool       `json:"ignore_not_found"`
	RetryTimeout   types.Time `json:"retry_timeout"`
	PollInterval   types.Time `json:"poll_interval"`
}

type webDriverConfig struct {
	Browser      string                 `json:"browser"`
	ManualStart  bool                   `json:"manual_start,omitempty"`
	BinaryPath   string                 `json:"binary_path,omitempty"`
	RemoteURL    string                 `json:"remote_url,omitempty"`
	Timeout      *types.Time            `json:"timeout,omitempty"`
	Capabalities map[string]interface{} `json:"capabilities,omitempty"`
}

type configParams struct {
	LogLevel            string           `json:"logging"`
	SoftAsserts         bool             `json:"soft_asserts"`
	ScreenshotDir       string           `json:"screenshot_dir,omitempty"`
	RaiseErrorsManually bool             `json:"raise_errors_automatically,omitempty"` //nolint:lll
	Runner              *runnerSettings  `json:"runner,omitempty"`
	Element             *elementSettings `json:"element,omitempty"`
	// TODO: Allow running multiple drivers.
	WebDriver *webDriverConfig `json:"webdriver,omitempty"`
}

var config *configParams

const defaultConfigPath = ".goseleniumrc.json"

func readConfig(configDirectory string) error {
	filePath := path.Join(configDirectory, defaultConfigPath)

	_, err := os.Stat(filePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "failed to stat config file")
		}

		logger.Info(
			"No config file found. Will create and use default config.",
		)

		c, err := createDefaultConfig()
		if err != nil {
			return errors.Wrap(err, "failed to create default config")
		}

		c.validateConfig()

		config = c

		return nil
	}

	if config == nil {
		c, err := readConfigFromFile(filePath)
		if err != nil {
			return errors.Wrap(err, "failed to read from config file")
		}

		c.validateConfig()

		config = c
	}

	return nil
}

func readConfigFromFile(filePath string) (*configParams, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return config, errors.Wrap(err, "failed to read config file")
	}

	var c configParams

	if err := json.Unmarshal(data, &c); err != nil {
		return config, errors.Wrap(err, "failed to parse config file")
	}

	return &c, nil
}

func createDefaultConfig() (*configParams, error) {
	c := &configParams{
		LogLevel:    logger.InfoLvl,
		SoftAsserts: false,
		WebDriver:   &webDriverConfig{Browser: "chrome"},
	}

	err := c.writeToConfig(defaultConfigPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to write config")
	}

	return c, nil
}

func (c *configParams) writeToConfig(configPath string) error {
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

func (c *configParams) validateConfig() {
	c.validateMain()
	c.validateRunner()
	c.validateElement()
	c.validateWebDriver()
}

func (c *configParams) validateMain() {
	if c.LogLevel == "" {
		logger.Warn(`"log_level" is not set. Defaulting to "info".`)

		c.LogLevel = "info"
	}
}

func (c *configParams) validateRunner() {
	if c.Runner == nil {
		c.Runner = &runnerSettings{ParallelRuns: 1}

		return
	}

	if c.Runner.ParallelRuns < 1 {
		logger.Warn(`"parallel_runs" is less than 1. Setting it to 1.`)

		c.Runner.ParallelRuns = 1
	}
}

func (c *configParams) validateElement() {
	defaultSettings := &elementSettings{
		IgnoreNotFound: false,
		SelectorType:   selector.CSS,
		RetryTimeout:   types.Time{Duration: 10 * time.Second},
		PollInterval:   types.Time{Duration: 500 * time.Millisecond},
	}

	if c.Element == nil {
		c.Element = &elementSettings{}
	}

	if c.Element.SelectorType == "" {
		logger.Warn(`"selector_type" is not set. Defaulting to "css".`)

		c.Element.SelectorType = defaultSettings.SelectorType
	}

	if c.Element.RetryTimeout.Duration == 0 {
		logger.Warn(`"retry_timeout" is not set. Defaulting to "10s".`)

		c.Element.RetryTimeout = defaultSettings.RetryTimeout
	}

	if c.Element.PollInterval.Duration == 0 {
		logger.Warn(`"poll_interval" is not set. Defaulting to "500ms".`)

		c.Element.PollInterval = defaultSettings.PollInterval
	}
}

func (c *configParams) validateWebDriver() {
	defaultSettings := &webDriverConfig{
		Browser:      "chrome",
		ManualStart:  false,
		BinaryPath:   "./chromedriver",
		RemoteURL:    "http://localhost:4444",
		Timeout:      &types.Time{Duration: 10 * time.Second},
		Capabalities: make(map[string]interface{}),
	}

	if c.WebDriver.BinaryPath == "" {
		logger.Warn(
			`"webdriver.binary" is not set. Defaulting to "chromedriver".`,
		)

		c.WebDriver.BinaryPath = defaultSettings.BinaryPath
	}

	if c.WebDriver.Timeout == nil {
		c.WebDriver.Timeout = &types.Time{}
	}

	if c.WebDriver.Timeout.Duration <= 0 {
		logger.Warn(`"timeout" is not set. Defaulting to "10s".`)

		c.WebDriver.Timeout = defaultSettings.Timeout
	}

	if c.WebDriver.RemoteURL == "" {
		logger.Warn(`"url" is not set. Defaulting to "http://localhost:4444".`)

		c.WebDriver.RemoteURL = defaultSettings.RemoteURL
	}
}
