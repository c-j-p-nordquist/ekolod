package probe

import (
	"net/http"
	"sync"
	"time"

	"github.com/c-j-p-nordquist/ekolod/pkg/metricspusher"
)

type LastRunChecker struct {
	lastRun time.Time
	mu      sync.RWMutex
}

func (c *LastRunChecker) Check() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return !c.lastRun.IsZero()
}

func (c *LastRunChecker) Name() string {
	return "last_run"
}

func (c *LastRunChecker) UpdateLastRun() {
	c.mu.Lock()
	c.lastRun = time.Now()
	c.mu.Unlock()
}

type CollectorReachableChecker struct{}

func (c *CollectorReachableChecker) Check() bool {
	resp, err := http.Get(metricspusher.GetCollectorURL() + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (c *CollectorReachableChecker) Name() string {
	return "collector_reachable"
}
