package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/c-j-p-nordquist/ekolod/internal/probe"
)

func MetricsHandler(probe probe.Probe) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		metrics := probe.GetMetrics()
		json.NewEncoder(w).Encode(metrics)
	}
}
