package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/c-j-p-nordquist/ekolod/pkg/probe"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/noop"
)

func TestHealthCheck(t *testing.T) {
	meter := noop.NewMeterProvider().Meter("test")
	prober := probe.NewProber(meter)
	exporter, _ := prometheus.New()
	s := NewServer(prober, exporter)

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handleHealth)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
