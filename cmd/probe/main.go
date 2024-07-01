package main

import (
	"log"
	"net/http"
	"os"

	"github.com/c-j-p-nordquist/ekolod/internal/handlers"
	"github.com/c-j-p-nordquist/ekolod/internal/probe"
	"github.com/c-j-p-nordquist/ekolod/pkg/config"
	"github.com/c-j-p-nordquist/ekolod/pkg/health"
	"github.com/c-j-p-nordquist/ekolod/pkg/logging"
	"github.com/c-j-p-nordquist/ekolod/pkg/metrics"
	"github.com/c-j-p-nordquist/ekolod/pkg/metricspusher"
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
	// Initialize health checker
	healthChecker := health.New()
	lastRunChecker := &probe.LastRunChecker{}
	healthChecker.AddChecker(lastRunChecker)
	healthChecker.AddChecker(&probe.CollectorReachableChecker{})

	collectorURL := os.Getenv("COLLECTOR_URL")
	if collectorURL == "" {
		log.Fatal("COLLECTOR_URL environment variable is not set")
	}

	if err := metricspusher.Init(collectorURL); err != nil {
		log.Fatalf("Failed to initialize metric pusher: %v", err)
	}

	// Start HTTP probe
	httpProbe := probe.NewHTTPProbe(targetPointers, lastRunChecker)

	// Run initial probe immediately
	httpProbe.RunProbe()

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
	mux.HandleFunc("/probe-metrics", handlers.ProbeMetricsHandler(httpProbe)) // JSON metrics
	mux.HandleFunc("/reload", handlers.ReloadHandler(httpProbe))
	mux.HandleFunc("/health", healthChecker.Handler())

	// Use CORS middleware
	corsMux := corsHandler(mux)

	// Start the server
	probePort := os.Getenv("PROBE_PORT")
	if probePort == "" {
		probePort = "8080"
	}
	logging.Info("Starting Probe HTTP server on :" + probePort)
	log.Fatal(http.ListenAndServe(":"+probePort, corsMux))
}
