package config

import (
	"io"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	LogLevel string   `yaml:"log_level"`
	Targets  []Target `yaml:"targets"`
}

type Target struct {
	Name              string        `yaml:"name"`
	URL               string        `yaml:"url"`
	Frequency         time.Duration `yaml:"frequency"`
	FailureTolerance  int           `yaml:"failure_tolerance"`
	RecoveryThreshold int           `yaml:"recovery_threshold"`
	Checks            []Check       `yaml:"checks"`
}

type Check struct {
	Path         string     `yaml:"path"`
	HTTPStatus   *Condition `yaml:"http_status,omitempty"`
	ResponseTime *Condition `yaml:"response_time,omitempty"`
	ResponseBody *Condition `yaml:"response_body,omitempty"`
}

type Condition struct {
	Type   string        `yaml:"condition"`
	Value  interface{}   `yaml:"value,omitempty"`
	Values []interface{} `yaml:"values,omitempty"`
}

type Threshold struct {
	Type  string      `yaml:"type"`
	Value interface{} `yaml:"value"`
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

	// Convert duration strings to time.Duration
	for i, target := range cfg.Targets {
		target.Frequency, err = time.ParseDuration(target.Frequency.String())
		if err != nil {
			return nil, err
		}
		for j, check := range target.Checks {
			if check.ResponseTime != nil && check.ResponseTime.Value != nil {
				if durationStr, ok := check.ResponseTime.Value.(string); ok {
					duration, err := time.ParseDuration(durationStr)
					if err != nil {
						return nil, err
					}
					cfg.Targets[i].Checks[j].ResponseTime.Value = duration
				}
			}
		}
	}

	return &cfg, nil
}
