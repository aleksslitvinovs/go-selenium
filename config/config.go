package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

//nolint:tagliatelle
type Config struct {
	Logging                      []string `json:"logging"`
	SoftAsserts                  bool     `json:"soft_asserts"`
	RaiseErrorsAutomaticatically bool     `json:"raise_errors_automaticially"`
}

const defaultFilename = "goseleniumrc.json"

// ReadConfig reads the config file and returns a Config struct.
func ReadConfig() (*Config, error) {
	_, err := os.Stat(defaultFilename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("No config file found. Using default config.")

			conf, err := createDefaultConfig()
			if err != nil {
				return nil, errors.Wrap(err, "failed to create default config")
			}

			return conf, nil
		}

		panic(err)
	}

	data, err := os.ReadFile(defaultFilename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}

	var c Config

	if err := json.Unmarshal(data, &c); err != nil {
		return nil, errors.Wrap(err, "failed to parse config file")
	}

	return &c, nil
}

func createDefaultConfig() (*Config, error) {
	f, err := os.Create(defaultFilename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create config file")
	}
	defer f.Close()

	c := Config{
		Logging:                      []string{"Asserts"},
		SoftAsserts:                  false,
		RaiseErrorsAutomaticatically: true,
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal config")
	}

	_, err = f.WriteString(string(data))
	if err != nil {
		return nil, errors.Wrap(err, "failed to write config")
	}

	return &c, nil
}
