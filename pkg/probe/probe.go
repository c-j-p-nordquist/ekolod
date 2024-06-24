package probe

import (
	"net/http"
	"time"
)

type Target struct {
	Name     string        `json:"name"`
	URL      string        `json:"url"`
	Interval time.Duration `json:"interval"`
	Timeout  time.Duration `json:"timeout"`
}

type ProbeResult struct {
	Target       string        `json:"target"`
	Status       string        `json:"status"`
	StatusCode   int           `json:"statusCode"`
	ResponseTime time.Duration `json:"responseTime"`
	Timestamp    time.Time     `json:"timestamp"`
}

func ProbeTarget(target Target) ProbeResult {
	client := &http.Client{
		Timeout: target.Timeout,
	}

	startTime := time.Now()
	resp, err := client.Get(target.URL)
	responseTime := time.Since(startTime)

	result := ProbeResult{
		Target:       target.Name,
		Timestamp:    time.Now(),
		ResponseTime: responseTime,
	}

	if err != nil {
		result.Status = "DOWN"
		result.StatusCode = 0
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		result.Status = "UP"
	} else {
		result.Status = "DOWN"
	}

	return result
}
