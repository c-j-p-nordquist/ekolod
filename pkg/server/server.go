package server

import (
	"encoding/json"
	"net/http"

	"github.com/c-j-p-nordquist/ekolod/pkg/probe"
)

type Server struct {
	prober *probe.Prober
}

func NewServer(prober *probe.Prober) *Server {
	return &Server{prober: prober}
}

func (s *Server) Start(addr string) error {
	http.HandleFunc("/results", s.handleResults)
	http.HandleFunc("/health", s.handleHealth)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleResults(w http.ResponseWriter, r *http.Request) {
	results := s.prober.GetResults()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
