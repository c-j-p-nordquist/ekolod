package probe

import "github.com/c-j-p-nordquist/ekolod/pkg/config"

type Probe interface {
	Start()
	Stop()
	UpdateTargets(targets []*config.Target)
	GetTargets() []config.Target
	GetMetrics() map[string]map[string]*ProbeResult
	AddTarget(target *config.Target)
	RemoveTarget(name string)
}

type ProbeResult struct {
	Duration float64
	Success  bool
	Message  string
}
