package probe

import (
	"github.com/c-j-p-nordquist/ekolod/pkg/config"
	"github.com/c-j-p-nordquist/ekolod/pkg/proberesult"
)

type Probe interface {
	Start()
	Stop()
	UpdateTargets(targets []*config.Target)
	GetTargets() []config.Target
	GetMetrics() map[string]map[string]*proberesult.ProbeResult
	AddTarget(target *config.Target)
	RemoveTarget(name string)
	RunProbe()
	UpdateTargetChecks(name string, checks []config.Check)
}
