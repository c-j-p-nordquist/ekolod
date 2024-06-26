package probe

import (
	"context"
	"net/http"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type Target struct {
	Name     string        `json:"name"`
	URL      string        `json:"url"`
	Interval time.Duration `json:"interval"`
	Timeout  time.Duration `json:"timeout"`
}

type Prober struct {
	targets    map[string]Target
	mu         sync.RWMutex
	meter      metric.Meter
	resultChan chan ProbeResult
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewProber(meter metric.Meter) *Prober {
	ctx, cancel := context.WithCancel(context.Background())
	return &Prober{
		targets:    make(map[string]Target),
		meter:      meter,
		resultChan: make(chan ProbeResult, 100),
		ctx:        ctx,
		cancel:     cancel,
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
}

func (p *Prober) RunProbes() {
	upGauge, _ := p.meter.Int64UpDownCounter("http_probe_up")
	responseTimeHistogram, _ := p.meter.Float64Histogram("http_probe_response_time")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			p.mu.RLock()
			targets := make([]Target, 0, len(p.targets))
			for _, target := range p.targets {
				targets = append(targets, target)
			}
			p.mu.RUnlock()

			for _, target := range targets {
				go func(t Target) {
					result := probeTarget(t)
					attrs := []attribute.KeyValue{
						attribute.String("target", t.Name),
						attribute.String("url", t.URL),
					}

					if result.Status == "UP" {
						upGauge.Add(p.ctx, 1, metric.WithAttributes(attrs...))
					} else {
						upGauge.Add(p.ctx, 0, metric.WithAttributes(attrs...))
					}

					responseTimeHistogram.Record(p.ctx, float64(result.ResponseTime.Milliseconds()), metric.WithAttributes(attrs...))

					select {
					case p.resultChan <- result:
					default:
						// If the channel is full, we'll skip this update
					}
				}(target)
			}
		}
	}
}

func (p *Prober) Stop() {
	p.cancel()
}

func (p *Prober) Results() <-chan ProbeResult {
	return p.resultChan
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

type ProbeResult struct {
	Target       string        `json:"target"`
	Status       string        `json:"status"`
	StatusCode   int           `json:"statusCode"`
	ResponseTime time.Duration `json:"responseTime"`
	Timestamp    time.Time     `json:"timestamp"`
}
