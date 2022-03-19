package config

type Config struct {
	Logging []string `json:"logging"`
}

func ReadConfig()
