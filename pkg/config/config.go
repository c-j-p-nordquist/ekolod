package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Targets []Target `yaml:"targets"`
	Server  Server   `yaml:"server"`
}

type Target struct {
	Name     string        `yaml:"name"`
	URL      string        `yaml:"url"`
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
}

type Server struct {
	Port int `yaml:"port"`
}

func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// Convert string durations to time.Duration
	for i, target := range config.Targets {
		interval, err := time.ParseDuration(target.Interval.String())
		if err != nil {
			return nil, fmt.Errorf("invalid interval for target %s: %w", target.Name, err)
		}
		config.Targets[i].Interval = interval

		timeout, err := time.ParseDuration(target.Timeout.String())
		if err != nil {
			return nil, fmt.Errorf("invalid timeout for target %s: %w", target.Name, err)
		}
		config.Targets[i].Timeout = timeout
	}

	return &config, nil
}
