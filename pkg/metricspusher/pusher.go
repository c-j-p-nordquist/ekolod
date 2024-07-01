package metricspusher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/c-j-p-nordquist/ekolod/pkg/config"
)

var collectorURL string

type ProbeResult struct {
	Duration       float64 `json:"duration"`
	Success        bool    `json:"success"`
	Message        string  `json:"message"`
	StatusCode     int     `json:"statusCode"`
	ContentLength  int64   `json:"contentLength"`
	TLSVersion     string  `json:"tlsVersion"`
	CertExpiryDays int     `json:"certExpiryDays"`
}

func Init(apiURL string) error {
	parsedURL, err := url.Parse(apiURL)
	if err != nil {
		return fmt.Errorf("invalid collector URL: %v", err)
	}
	collectorURL = parsedURL.String()
	return nil
}

func GetCollectorURL() string {
	return collectorURL
}

func PushMetricsToCollector(target *config.Target, check config.Check, result *ProbeResult) error {
	if collectorURL == "" {
		return fmt.Errorf("collector URL is not set")
	}

	payload := map[string]interface{}{
		"target": target.Name,
		"check":  check.Path,
		"result": result,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Post(collectorURL+"/metrics", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to push metrics: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("collector responded with status code: %d", resp.StatusCode)
	}

	return nil
}
