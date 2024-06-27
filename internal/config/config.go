package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	LogLevel string   `yaml:"log_level"`
	Targets  []string `yaml:"targets"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
