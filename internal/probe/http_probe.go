package probe

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/c-j-p-nordquist/ekolod/pkg/checker"
	"github.com/c-j-p-nordquist/ekolod/pkg/config"
	"github.com/c-j-p-nordquist/ekolod/pkg/logging"
	"github.com/c-j-p-nordquist/ekolod/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type HTTPProbe struct {
	targets      []*config.Target
	mu           sync.Mutex
	metrics      map[string]map[string]*ProbeResult
	stopChannels map[string]chan struct{}
}

func NewHTTPProbe(targets []*config.Target) Probe {
	probe := &HTTPProbe{
		targets:      targets,
		metrics:      make(map[string]map[string]*ProbeResult),
		stopChannels: make(map[string]chan struct{}),
	}
	probe.Start()
	return probe
}

func (p *HTTPProbe) Start() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, target := range p.targets {
		if _, exists := p.stopChannels[target.Name]; !exists {
			stopChan := make(chan struct{})
			p.stopChannels[target.Name] = stopChan
			go p.probeTarget(target, stopChan)
		}
	}
}

func (p *HTTPProbe) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, stopChan := range p.stopChannels {
		close(stopChan)
	}
	p.stopChannels = make(map[string]chan struct{})
}

func (p *HTTPProbe) UpdateTargets(targets []*config.Target) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Stop probing old targets that are not in the new configuration
	for name, stopChan := range p.stopChannels {
		found := false
		for _, target := range targets {
			if target.Name == name {
				found = true
				break
			}
		}
		if !found {
			close(stopChan)
			delete(p.stopChannels, name)
			delete(p.metrics, name)
		}
	}

	// Start probing new targets
	for _, target := range targets {
		if _, exists := p.stopChannels[target.Name]; !exists {
			stopChan := make(chan struct{})
			p.stopChannels[target.Name] = stopChan
			go p.probeTarget(target, stopChan)
		}
	}

	p.targets = targets
}

func (p *HTTPProbe) GetTargets() []config.Target {
	p.mu.Lock()
	defer p.mu.Unlock()
	targets := make([]config.Target, len(p.targets))
	for i, t := range p.targets {
		targets[i] = *t
	}
	return targets
}

func (p *HTTPProbe) AddTarget(target *config.Target) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.targets = append(p.targets, target)
	if _, exists := p.stopChannels[target.Name]; !exists {
		stopChan := make(chan struct{})
		p.stopChannels[target.Name] = stopChan
		go p.probeTarget(target, stopChan)
	}
}

func (p *HTTPProbe) RemoveTarget(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i, t := range p.targets {
		if t.Name == name {
			p.targets = append(p.targets[:i], p.targets[i+1:]...)
			if stopChan, exists := p.stopChannels[name]; exists {
				close(stopChan)
				delete(p.stopChannels, name)
			}
			delete(p.metrics, name)
			break
		}
	}
}

func (p *HTTPProbe) GetMetrics() map[string]map[string]*ProbeResult {
	p.mu.Lock()
	defer p.mu.Unlock()
	copy := make(map[string]map[string]*ProbeResult)
	for k, v := range p.metrics {
		copy[k] = make(map[string]*ProbeResult)
		for kk, vv := range v {
			copy[k][kk] = vv
		}
	}
	return copy
}

func (p *HTTPProbe) probeTarget(target *config.Target, stopChan chan struct{}) {
	ticker := time.NewTicker(target.Frequency)
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			for _, check := range target.Checks {
				result := p.runCheck(target, convertConfigCheckToCheckerCheck(check))

				p.mu.Lock()
				if p.metrics[target.Name] == nil {
					p.metrics[target.Name] = make(map[string]*ProbeResult)
				}
				p.metrics[target.Name][check.Path] = result
				p.mu.Unlock()

				metrics.HttpRequestDuration.With(prometheus.Labels{
					"target": target.Name,
					"path":   check.Path,
					"method": "GET", // Assuming GET for now, update if you add support for other methods
				}).Observe(result.Duration)

				if !result.Success {
					logging.Warn(fmt.Sprintf("Check failed for target '%s', path '%s': %s", target.Name, check.Path, result.Message))
				}
			}
		}
	}
}

func (p *HTTPProbe) runCheck(target *config.Target, check checker.Check) *ProbeResult {
	url := target.URL + check.Path
	start := time.Now()
	resp, err := http.Get(url)
	duration := time.Since(start)

	if err != nil {
		return &ProbeResult{
			Duration: duration.Seconds(),
			Success:  false,
			Message:  fmt.Sprintf("HTTP request failed: %v", err),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ProbeResult{
			Duration: duration.Seconds(),
			Success:  false,
			Message:  fmt.Sprintf("Failed to read response body: %v", err),
		}
	}

	checkerResponse := checker.Response{
		StatusCode: resp.StatusCode,
		Body:       string(body),
		Duration:   duration,
	}

	result := checker.EvaluateCheck(check, checkerResponse)

	return &ProbeResult{
		Duration: duration.Seconds(),
		Success:  result.Success,
		Message:  result.Message,
	}
}

func convertConfigCheckToCheckerCheck(configCheck config.Check) checker.Check {
	return checker.Check{
		Path:         configCheck.Path,
		HTTPStatus:   convertCondition(configCheck.HTTPStatus),
		ResponseTime: convertCondition(configCheck.ResponseTime),
		ResponseBody: convertCondition(configCheck.ResponseBody),
	}
}

func convertCondition(configCondition *config.Condition) *checker.Condition {
	if configCondition == nil {
		return nil
	}
	return &checker.Condition{
		Type:   configCondition.Type,
		Value:  configCondition.Value,
		Values: configCondition.Values,
	}
}
