package handlers

import (
	"net/http"
	"sync"

	"github.com/c-j-p-nordquist/ekolod/internal/config"
	"github.com/c-j-p-nordquist/ekolod/internal/logging"
	"github.com/c-j-p-nordquist/ekolod/internal/probe"
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
	targetPointers := make([]*config.Target, len(cfg.Targets))
	for i := range cfg.Targets {
		targetPointers[i] = &cfg.Targets[i]
	}
	probe.UpdateTargets(targetPointers)
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

		targetPointers := make([]*config.Target, len(cfg.Targets))
		for i := range cfg.Targets {
			targetPointers[i] = &cfg.Targets[i]
		}
		probe.UpdateTargets(targetPointers)

		logging.Info("Configuration reloaded successfully")
		w.Write([]byte("Configuration reloaded successfully"))
	}
}
