package probe

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/c-j-p-nordquist/ekolod/internal/logging"
	"github.com/c-j-p-nordquist/ekolod/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type HTTPProbe struct {
	targets []string
	mu      sync.Mutex
	metrics map[string]float64
}

func NewHTTPProbe(targets []string) *HTTPProbe {
	return &HTTPProbe{targets: targets, metrics: make(map[string]float64)}
}

func (p *HTTPProbe) Start() {
	for _, target := range p.targets {
		go p.probeTarget(target)
	}
}

func (p *HTTPProbe) GetTargets() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]string{}, p.targets...) // Return a copy of the targets slice
}

func (p *HTTPProbe) AddTarget(target string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.targets = append(p.targets, target)
	go p.probeTarget(target)
}

func (p *HTTPProbe) RemoveTarget(target string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, t := range p.targets {
		if t == target {
			p.targets = append(p.targets[:i], p.targets[i+1:]...)
			break
		}
	}
}

func (p *HTTPProbe) UpdateTargets(targets []string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.targets = targets
	for _, target := range p.targets {
		go p.probeTarget(target)
	}
}

func (p *HTTPProbe) probeTarget(target string) {
	for {
		start := time.Now()
		resp, err := p.httpGetWithRetry(target, 3, 2*time.Second)
		duration := time.Since(start).Seconds()

		if err == nil && resp != nil {
			resp.Body.Close()
		}

		p.mu.Lock()
		p.metrics[target] = duration
		p.mu.Unlock()

		metrics.HttpRequestDuration.With(prometheus.Labels{"method": "GET", "endpoint": target}).Observe(duration)
		time.Sleep(10 * time.Second)
	}
}

func (p *HTTPProbe) httpGetWithRetry(url string, retries int, delay time.Duration) (*http.Response, error) {
	for i := 0; i < retries; i++ {
		resp, err := http.Get(url)
		if err == nil {
			return resp, nil
		}
		logging.Warn("Retrying " + url + " due to error: " + err.Error())
		time.Sleep(delay)
	}
	err := fmt.Errorf("failed to GET %s after %d attempts", url, retries)
	logging.Error(err)
	return nil, err
}

func (p *HTTPProbe) GetMetrics() map[string]float64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	copy := make(map[string]float64)
	for k, v := range p.metrics {
		copy[k] = v
	}
	return copy
}
