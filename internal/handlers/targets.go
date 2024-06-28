package handlers

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/c-j-p-nordquist/ekolod/internal/probe"
	"github.com/c-j-p-nordquist/ekolod/pkg/config"
	"github.com/c-j-p-nordquist/ekolod/pkg/logging"
)

var (
	targetList []config.Target
	targetMu   sync.Mutex
)

func InitTargets(probe probe.Probe) {
	targetList = probe.GetTargets()
}

func ListTargetsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetMu.Lock()
		defer targetMu.Unlock()

		json.NewEncoder(w).Encode(targetList)
	}
}

func AddTargetHandler(probe probe.Probe) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var target config.Target
		err := json.NewDecoder(r.Body).Decode(&target)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		targetMu.Lock()
		defer targetMu.Unlock()

		probe.AddTarget(&target)
		targetList = probe.GetTargets() // Update targetList with the new targets

		// Update the global cfg variable
		mu.Lock()
		cfg.Targets = append(cfg.Targets, target)
		mu.Unlock()

		logging.Info("Added new target: " + target.Name)

		w.WriteHeader(http.StatusCreated)
	}
}

func RemoveTargetHandler(probe probe.Probe) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var target struct {
			Name string `json:"name"`
		}
		err := json.NewDecoder(r.Body).Decode(&target)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		targetMu.Lock()
		defer targetMu.Unlock()

		probe.RemoveTarget(target.Name)
		targetList = probe.GetTargets() // Update targetList with the new targets

		// Update the global cfg variable
		mu.Lock()
		for i, t := range cfg.Targets {
			if t.Name == target.Name {
				cfg.Targets = append(cfg.Targets[:i], cfg.Targets[i+1:]...)
				break
			}
		}
		mu.Unlock()

		logging.Info("Removed target: " + target.Name)
		w.WriteHeader(http.StatusOK)
	}
}
