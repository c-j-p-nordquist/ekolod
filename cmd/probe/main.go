package main

import (
	"log"
	"net/http"

	"github.com/c-j-p-nordquist/ekolod/internal/config"
	"github.com/c-j-p-nordquist/ekolod/internal/handlers"
	"github.com/c-j-p-nordquist/ekolod/internal/logging"
	"github.com/c-j-p-nordquist/ekolod/internal/metrics"
	"github.com/c-j-p-nordquist/ekolod/internal/probe"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		logging.Error(err)
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize logging
	logging.InitLogger(cfg.LogLevel)

	// Initialize metrics
	metrics.InitMetrics()

	// Convert cfg.Targets to []*config.Target
	targetPointers := make([]*config.Target, len(cfg.Targets))
	for i := range cfg.Targets {
		targetPointers[i] = &cfg.Targets[i]
	}

	// Start HTTP probe
	httpProbe := probe.NewHTTPProbe(targetPointers)

	// Initialize handlers
	handlers.Init("configs/config.yaml", httpProbe)

	// Setup CORS middleware
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	// Create a new ServeMux
	mux := http.NewServeMux()
	mux.Handle("/metrics", metrics.Handler())
	mux.HandleFunc("/metrics/data", handlers.MetricsHandler(httpProbe))
	mux.HandleFunc("/reload", handlers.ReloadHandler(httpProbe))

	// Use CORS middleware
	corsMux := corsHandler(mux)

	logging.Info("Starting HTTP server for metrics and API endpoints...")
	log.Fatal(http.ListenAndServe(":8080", corsMux))
}
