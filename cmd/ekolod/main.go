package main

import (
	"context"
	"log"
	"time"

	"github.com/c-j-p-nordquist/ekolod/pkg/probe"
	"github.com/c-j-p-nordquist/ekolod/pkg/server"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

func main() {
	ctx := context.Background()

	exporter, err := prometheus.New()
	if err != nil {
		log.Fatalf("Failed to create Prometheus exporter: %v", err)
	}

	meterProvider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := meterProvider.Meter("ekolod")

	prober := probe.NewProber(meter)

	targets := []probe.Target{
		{
			Name:     "Example",
			URL:      "https://example.com",
			Interval: 10 * time.Second,
			Timeout:  5 * time.Second,
		},
		{
			Name:     "Google",
			URL:      "https://www.google.com",
			Interval: 15 * time.Second,
			Timeout:  5 * time.Second,
		},
	}

	for _, target := range targets {
		prober.AddTarget(target)
	}

	go prober.RunProbes(ctx)

	srv := server.NewServer(prober, exporter)

	log.Println("Starting server on :8080")
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
