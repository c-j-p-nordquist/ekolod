package handlers

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/c-j-p-nordquist/ekolod/internal/logging"
	"github.com/c-j-p-nordquist/ekolod/internal/probe"
)

var (
	targetList []string
	targetMu   sync.Mutex
)

func InitTargets(probe *probe.HTTPProbe) {
	targetList = probe.GetTargets()
}

func ListTargetsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetMu.Lock()
		defer targetMu.Unlock()

		json.NewEncoder(w).Encode(targetList)
	}
}

func AddTargetHandler(probe *probe.HTTPProbe) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var target struct {
			URL string `json:"url"`
		}
		err := json.NewDecoder(r.Body).Decode(&target)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		targetMu.Lock()
		defer targetMu.Unlock()

		probe.AddTarget(target.URL)
		targetList = probe.GetTargets() // Update targetList with the new targets
		logging.Info("Added new target: " + target.URL)

		w.WriteHeader(http.StatusCreated)
	}
}

func RemoveTargetHandler(probe *probe.HTTPProbe) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var target struct {
			URL string `json:"url"`
		}
		err := json.NewDecoder(r.Body).Decode(&target)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		targetMu.Lock()
		defer targetMu.Unlock()

		probe.RemoveTarget(target.URL)
		targetList = probe.GetTargets() // Update targetList with the new targets

		logging.Info("Removed target: " + target.URL)
		w.WriteHeader(http.StatusOK)
	}
}
