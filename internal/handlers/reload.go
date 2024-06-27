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

func Init(cfgPath string, probe *probe.HTTPProbe) {
	var err error
	cfg, err = config.LoadConfig(cfgPath)
	if err != nil {
		logging.Error(err)
		return
	}
	// Initial probe setup
	probe.UpdateTargets(cfg.Targets)
}

func ReloadHandler(probe *probe.HTTPProbe) http.HandlerFunc {
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

		probe.UpdateTargets(cfg.Targets)
		logging.Info("Configuration reloaded successfully")
		w.Write([]byte("Configuration reloaded successfully"))
	}
}
