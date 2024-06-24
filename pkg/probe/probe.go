package probe

import (
	"net/http"
	"sync"
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

type Prober struct {
	targets map[string]Target
	results map[string]ProbeResult
	mu      sync.RWMutex
}

func NewProber() *Prober {
	return &Prober{
		targets: make(map[string]Target),
		results: make(map[string]ProbeResult),
	}
}

func (p *Prober) AddTarget(target Target) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.targets[target.Name] = target
}

func (p *Prober) RemoveTarget(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.targets, name)
	delete(p.results, name)
}

func (p *Prober) GetResults() map[string]ProbeResult {
	p.mu.RLock()
	defer p.mu.RUnlock()
	results := make(map[string]ProbeResult)
	for k, v := range p.results {
		results[k] = v
	}
	return results
}

func (p *Prober) RunProbes() {
	for {
		p.mu.RLock()
		targets := make([]Target, 0, len(p.targets))
		for _, target := range p.targets {
			targets = append(targets, target)
		}
		p.mu.RUnlock()

		for _, target := range targets {
			go func(t Target) {
				result := probeTarget(t)
				p.mu.Lock()
				p.results[t.Name] = result
				p.mu.Unlock()
			}(target)
		}

		time.Sleep(1 * time.Second) // Wait before next round of probes
	}
}

func probeTarget(target Target) ProbeResult {
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
