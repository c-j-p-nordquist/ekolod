package handlers

import (
	"net/http"
	"sync"

	"github.com/c-j-p-nordquist/ekolod/internal/probe"
	"github.com/c-j-p-nordquist/ekolod/pkg/config"
	"github.com/c-j-p-nordquist/ekolod/pkg/logging"
)

var (
	cfg *config.Config
	mu  sync.Mutex
)

func Init(cfgPath string, probe probe.Probe) {
	var err error
	cfg, err = config.LoadConfig(cfgPath)
	if err != nil {
		logging.Error(err)
		return
	}
	updateProbeAndTargetList(probe)
}

func updateProbeAndTargetList(probe probe.Probe) {
	newTargets := make([]*config.Target, len(cfg.Targets))
	for i := range cfg.Targets {
		newTargets[i] = &cfg.Targets[i]
	}

	oldTargets := probe.GetTargets()

	// Update existing targets and add new ones
	for _, newTarget := range newTargets {
		found := false
		for _, oldTarget := range oldTargets {
			if oldTarget.Name == newTarget.Name {
				probe.UpdateTargetChecks(newTarget.Name, newTarget.Checks)
				found = true
				break
			}
		}
		if !found {
			probe.AddTarget(newTarget)
		}
	}

	// Remove targets that no longer exist
	for _, oldTarget := range oldTargets {
		found := false
		for _, newTarget := range newTargets {
			if oldTarget.Name == newTarget.Name {
				found = true
				break
			}
		}
		if !found {
			probe.RemoveTarget(oldTarget.Name)
		}
	}

	// Update targetList
	targetMu.Lock()
	targetList = probe.GetTargets()
	targetMu.Unlock()
}

func ReloadHandler(probe probe.Probe) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		var err error
		cfg, err = config.LoadConfig("configs/config.yaml")
		if err != nil {
			logging.Error(err)
			http.Error(w, "Failed to reload config", http.StatusInternalServerError)
			return
		}

		updateProbeAndTargetList(probe)

		logging.Info("Configuration reloaded successfully")
		w.Write([]byte("Configuration reloaded successfully"))
	}
}
