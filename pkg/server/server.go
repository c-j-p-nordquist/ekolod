package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/c-j-p-nordquist/ekolod/pkg/probe"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
)

type Server struct {
	prober   *probe.Prober
	exporter *prometheus.Exporter
	reloadCh chan struct{}
	mu       sync.Mutex
}

func NewServer(prober *probe.Prober, exporter *prometheus.Exporter) *Server {
	return &Server{
		prober:   prober,
		exporter: exporter,
		reloadCh: make(chan struct{}, 1),
	}
}

func (s *Server) Start(addr string) error {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", s.handleHealth)
	http.HandleFunc("/reload", s.handleReload)

	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	log.Printf("Health check requested from %s", r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) handleReload(w http.ResponseWriter, r *http.Request) {
	log.Printf("Reload requested from %s", r.RemoteAddr)

	if r.Method != http.MethodPost {
		log.Printf("Invalid method for reload: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.mu.Lock()
	select {
	case s.reloadCh <- struct{}{}:
		log.Println("Reload triggered successfully")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Reload triggered"))
	default:
		log.Println("Reload already in progress, request ignored")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("Reload already in progress"))
	}
	s.mu.Unlock()
}

func (s *Server) WaitForReload() <-chan struct{} {
	return s.reloadCh
}

func (s *Server) UpdateProber(prober *probe.Prober) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prober = prober
	log.Println("Prober updated with new configuration")
}
