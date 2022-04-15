package config

import (
	"time"

	"github.com/theRealAlpaca/go-selenium/logger"
	"github.com/theRealAlpaca/go-selenium/selector"
	"github.com/theRealAlpaca/go-selenium/types"
)

func (c *config) validateConfig() {
	c.validateMain()
	c.validateRunner()
	c.validateElement()
	c.validateWebDriver()
}

func (c *config) validateMain() {
	if c.LogLevel == "" {
		logger.Warn(`"log_level" is not set. Defaulting to "info".`)

		c.LogLevel = "info"
	}
}

func (c *config) validateRunner() {
	if c.Runner == nil {
		c.Runner = &RunnerSettings{ParallelRuns: 1}

		return
	}

	if c.Runner.ParallelRuns < 1 {
		logger.Warn(`"parallel_runs" is less than 1. Setting it to 1.`)

		c.Runner.ParallelRuns = 1
	}
}

func (c *config) validateElement() {
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
}

func (c *config) validateWebDriver() {
	if c.WebDriver.PathToBinary == "" {
		logger.Warn(
			`"webdriver.binary" is not set. Defaulting to "chromedriver".`,
		)

		c.WebDriver.PathToBinary = "chromedriver"
	}

	if c.WebDriver.Timeout.Duration <= 0 {
		logger.Warn(`"timeout" is not set. Defaulting to "10s".`)

		c.WebDriver.Timeout = types.Time{Duration: 10 * time.Second}
	}

	if c.WebDriver.URL == "" {
		logger.Warn(`"url" is not set. Defaulting to "http://localhost:4444".`)

		c.WebDriver.URL = "http://localhost:4444"
	}
}
