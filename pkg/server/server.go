package server

import (
	"net/http"

	"github.com/c-j-p-nordquist/ekolod/pkg/probe"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
)

type Server struct {
	prober   *probe.Prober
	exporter *prometheus.Exporter
}

func NewServer(prober *probe.Prober, exporter *prometheus.Exporter) *Server {
	return &Server{
		prober:   prober,
		exporter: exporter,
	}
}

func (s *Server) Start(addr string) error {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", s.handleHealth)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
