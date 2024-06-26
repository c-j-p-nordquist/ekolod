package probe

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.opentelemetry.io/otel/metric/noop"
)

func TestProbeTarget(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	target := Target{
		Name:     "Test",
		URL:      ts.URL,
		Interval: 1 * time.Second,
		Timeout:  500 * time.Millisecond,
	}

	result := probeTarget(target)

	if result.Status != "UP" {
		t.Errorf("Expected status UP, got %s", result.Status)
	}

	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, result.StatusCode)
	}
}

func TestProber(t *testing.T) {
	meter := noop.NewMeterProvider().Meter("test")
	prober := NewProber(meter)

	target := Target{
		Name:     "Test",
		URL:      "http://example.com",
		Interval: 1 * time.Second,
		Timeout:  500 * time.Millisecond,
	}

	prober.AddTarget(target)

	if len(prober.targets) != 1 {
		t.Errorf("Expected 1 target, got %d", len(prober.targets))
	}

	prober.RemoveTarget("Test")

	if len(prober.targets) != 0 {
		t.Errorf("Expected 0 targets, got %d", len(prober.targets))
	}

	// Test RunProbes and Stop
	prober.AddTarget(target)
	go prober.RunProbes()

	// Wait for at least one probe cycle
	time.Sleep(1500 * time.Millisecond)

	// Check if we're receiving results
	select {
	case result := <-prober.Results():
		if result.Target != "Test" {
			t.Errorf("Expected target Test, got %s", result.Target)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timed out waiting for probe result")
	}

	// Test Stop
	prober.Stop()

	// Ensure no more results are coming after stop
	select {
	case <-prober.Results():
		t.Error("Received result after stopping prober")
	case <-time.After(1500 * time.Millisecond):
		// This is the expected behavior
	}
}

func TestProberResults(t *testing.T) {
	meter := noop.NewMeterProvider().Meter("test")
	prober := NewProber(meter)

	resultChan := prober.Results()

	if resultChan == nil {
		t.Error("Results() returned nil channel")
	}
}
