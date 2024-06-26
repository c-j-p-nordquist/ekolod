package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/c-j-p-nordquist/ekolod/pkg/probe"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
)

type Server struct {
	prober     *probe.Prober
	exporter   *prometheus.Exporter
	reloadCh   chan struct{}
	mu         sync.Mutex
	sseClients map[chan probe.ProbeResult]bool
	sseMu      sync.RWMutex
}

func NewServer(prober *probe.Prober, exporter *prometheus.Exporter) *Server {
	s := &Server{
		prober:     prober,
		exporter:   exporter,
		reloadCh:   make(chan struct{}, 1),
		sseClients: make(map[chan probe.ProbeResult]bool),
	}
	go s.broadcastResults()
	return s
}

func (s *Server) broadcastResults() {
	for result := range s.prober.Results() {
		s.BroadcastProbeResult(result)
	}
}

func (s *Server) Start(addr string) error {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", s.handleHealth)
	http.HandleFunc("/reload", s.handleReload)
	http.HandleFunc("/events", s.handleSSE)

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

func (s *Server) handleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	resultChan := make(chan probe.ProbeResult, 10)

	s.sseMu.Lock()
	s.sseClients[resultChan] = true
	s.sseMu.Unlock()

	defer func() {
		s.sseMu.Lock()
		delete(s.sseClients, resultChan)
		s.sseMu.Unlock()
		close(resultChan)
	}()

	for {
		select {
		case <-r.Context().Done():
			return
		case result := <-resultChan:
			data, err := json.Marshal(result)
			if err != nil {
				log.Printf("Error marshaling probe result: %v", err)
				continue
			}
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}

func (s *Server) BroadcastProbeResult(result probe.ProbeResult) {
	s.sseMu.RLock()
	defer s.sseMu.RUnlock()

	for client := range s.sseClients {
		select {
		case client <- result:
		default:
			// If the client's channel is full, we'll skip this update for this client
		}
	}
}

func (s *Server) WaitForReload() <-chan struct{} {
	return s.reloadCh
}

func (s *Server) UpdateProber(newProber *probe.Prober) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Stop the old prober
	s.prober.Stop()

	// Update to the new prober
	s.prober = newProber

	// Start a new broadcastResults goroutine
	go s.broadcastResults()

	log.Println("Prober updated with new configuration")
}
