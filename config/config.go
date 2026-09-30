package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const defaultPath = "config.yaml"

type Config struct {
	AppName     string
	APIPort     int
	Environment string
}

type fileConfig struct {
	AppName     string `yaml:"app_name"`
	APIPort     int    `yaml:"apiport"`
	Environment string `yaml:"environment"`
}

func Load(path string) (Config, error) {
	if path == "" {
		path = defaultPath
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var file fileConfig
	if err := yaml.Unmarshal(data, &file); err != nil {
		return Config{}, err
	}

	if file.AppName == "" {
		file.AppName = "auto-workflow"
	}
	if file.APIPort == 0 {
		file.APIPort = 8080
	}
	if file.Environment == "" {
		file.Environment = "development"
	}

	return Config{
		AppName:     file.AppName,
		APIPort:     file.APIPort,
		Environment: file.Environment,
	}, nil
}
