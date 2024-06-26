package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func TestSSEEndpoint(t *testing.T) {
	meter := noop.NewMeterProvider().Meter("test")
	prober := probe.NewProber(meter)
	exporter, _ := prometheus.New()
	s := NewServer(prober, exporter)

	// Add a test target
	prober.AddTarget(probe.Target{
		Name: "Test",
		URL:  "http://example.com",
	})

	// Start the prober
	go prober.RunProbes()

	req, err := http.NewRequest("GET", "/events", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handleSSE)

	// Use a context with timeout to automatically end the test
	ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	go handler.ServeHTTP(rr, req)

	// Wait for some data to be written
	time.Sleep(1500 * time.Millisecond)

	// Check the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check if the response contains SSE data
	if !strings.Contains(rr.Body.String(), "data:") {
		t.Errorf("SSE data not found in response")
	}

	// Attempt to parse the data as JSON
	dataLines := strings.Split(strings.TrimSpace(rr.Body.String()), "\n\n")
	for _, line := range dataLines {
		if strings.HasPrefix(line, "data:") {
			jsonData := strings.TrimPrefix(line, "data:")
			var result probe.ProbeResult
			err := json.Unmarshal([]byte(jsonData), &result)
			if err != nil {
				t.Errorf("Failed to parse SSE data as JSON: %v", err)
			}
			// Check if the parsed data has the expected structure
			if result.Target != "Test" {
				t.Errorf("Unexpected target in SSE data: got %v, want %v", result.Target, "Test")
			}
		}
	}
}

func TestUpdateProber(t *testing.T) {
	meter := noop.NewMeterProvider().Meter("test")
	oldProber := probe.NewProber(meter)
	exporter, _ := prometheus.New()
	s := NewServer(oldProber, exporter)

	newProber := probe.NewProber(meter)
	s.UpdateProber(newProber)

	if s.prober != newProber {
		t.Errorf("Prober was not updated")
	}
}
