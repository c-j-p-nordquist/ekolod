package probe

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/c-j-p-nordquist/ekolod/pkg/checkconverter"
	"github.com/c-j-p-nordquist/ekolod/pkg/checker"
	"github.com/c-j-p-nordquist/ekolod/pkg/config"
	"github.com/c-j-p-nordquist/ekolod/pkg/httputils"
	"github.com/c-j-p-nordquist/ekolod/pkg/logging"
	"github.com/c-j-p-nordquist/ekolod/pkg/metrics"
	"github.com/c-j-p-nordquist/ekolod/pkg/metricspusher"
	"github.com/c-j-p-nordquist/ekolod/pkg/proberesult"
	"github.com/c-j-p-nordquist/ekolod/pkg/tlsutils"
)

type HTTPProbe struct {
	targets        []*config.Target
	mu             sync.Mutex
	metrics        map[string]map[string]*proberesult.ProbeResult
	stopChannels   map[string]chan struct{}
	lastRunChecker *LastRunChecker
}

func NewHTTPProbe(targets []*config.Target, lastRunChecker *LastRunChecker) *HTTPProbe {
	probe := &HTTPProbe{
		targets:        targets,
		metrics:        make(map[string]map[string]*proberesult.ProbeResult),
		stopChannels:   make(map[string]chan struct{}),
		lastRunChecker: lastRunChecker,
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

func (p *HTTPProbe) GetMetrics() map[string]map[string]*proberesult.ProbeResult {
	p.mu.Lock()
	defer p.mu.Unlock()
	copy := make(map[string]map[string]*proberesult.ProbeResult)
	for k, v := range p.metrics {
		copy[k] = make(map[string]*proberesult.ProbeResult)
		for kk, vv := range v {
			copy[k][kk] = vv
		}
	}
	return copy
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

func (p *HTTPProbe) RunProbe() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, target := range p.targets {
		for _, check := range target.Checks {
			result := p.runCheck(target, checkconverter.ConvertConfigCheckToCheckerCheck(check))

			if p.metrics[target.Name] == nil {
				p.metrics[target.Name] = make(map[string]*proberesult.ProbeResult)
			}
			p.metrics[target.Name][check.Path] = result

			metrics.UpdatePrometheusMetrics(target, check, result)

			if !result.Success {
				logging.Warn(fmt.Sprintf("Check failed for target '%s', path '%s': %s", target.Name, check.Path, result.Message))
			}
		}
	}
}

func (p *HTTPProbe) UpdateTargetChecks(name string, checks []config.Check) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i, target := range p.targets {
		if target.Name == name {
			p.targets[i].Checks = checks
			if stopChan, exists := p.stopChannels[name]; exists {
				close(stopChan)
				delete(p.stopChannels, name)
			}
			stopChan := make(chan struct{})
			p.stopChannels[name] = stopChan
			go p.probeTarget(p.targets[i], stopChan)
			break
		}
	}
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
				result := p.runCheck(target, checkconverter.ConvertConfigCheckToCheckerCheck(check))

				p.mu.Lock()
				if p.metrics[target.Name] == nil {
					p.metrics[target.Name] = make(map[string]*proberesult.ProbeResult)
				}
				p.metrics[target.Name][check.Path] = result
				p.mu.Unlock()

				metrics.UpdatePrometheusMetrics(target, check, result)

				pusherResult := &metricspusher.ProbeResult{
					Duration:       result.Duration,
					Success:        result.Success,
					Message:        result.Message,
					StatusCode:     result.StatusCode,
					ContentLength:  result.ContentLength,
					TLSVersion:     result.TLSVersion,
					CertExpiryDays: result.CertExpiryDays,
				}

				err := metricspusher.PushMetricsToCollector(target, check, pusherResult)
				if err != nil {
					logging.Error(fmt.Errorf("failed to push metrics for target '%s', path '%s': %v", target.Name, check.Path, err))
				}

				if !result.Success {
					logging.Warn(fmt.Sprintf("Check failed for target '%s', path '%s': %s", target.Name, check.Path, result.Message))
				}
			}
		}
	}
}

func (p *HTTPProbe) runCheck(target *config.Target, check checker.Check) *proberesult.ProbeResult {
	url := target.URL + check.Path
	start := time.Now()

	client := httputils.NewHTTPClient()

	resp, err := client.Get(url)
	duration := time.Since(start)

	result := proberesult.New(duration.Seconds())

	if err != nil {
		result.SetMessage(fmt.Sprintf("HTTP request failed: %v", err))
		return result
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.SetMessage(fmt.Sprintf("Failed to read response body: %v", err))
		return result
	}

	result.SetStatusCode(resp.StatusCode)
	result.SetContentLength(resp.ContentLength)

	if resp.TLS != nil {
		result.SetTLSVersion(tlsutils.VersionToString(resp.TLS.Version))
		if len(resp.TLS.PeerCertificates) > 0 {
			result.SetCertExpiryDays(tlsutils.CalculateCertExpiryDays(resp.TLS.PeerCertificates[0].NotAfter))
		}
	}

	checkerResponse := checker.Response{
		StatusCode: resp.StatusCode,
		Body:       string(body),
		Duration:   duration,
	}

	checkResult := checker.EvaluateCheck(check, checkerResponse)
	result.SetSuccess(checkResult.Success)
	result.SetMessage(checkResult.Message)

	p.lastRunChecker.UpdateLastRun()

	return result
}
