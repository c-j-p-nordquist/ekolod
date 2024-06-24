package main

import (
	"context"
	"fmt"
	"log"

	"github.com/c-j-p-nordquist/ekolod/pkg/config"
	"github.com/c-j-p-nordquist/ekolod/pkg/probe"
	"github.com/c-j-p-nordquist/ekolod/pkg/server"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx := context.Background()

	exporter, err := prometheus.New()
	if err != nil {
		log.Fatalf("Failed to create Prometheus exporter: %v", err)
	}

	meterProvider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := meterProvider.Meter("ekolod")

	prober := probe.NewProber(meter)

	for _, target := range cfg.Targets {
		prober.AddTarget(probe.Target{
			Name:     target.Name,
			URL:      target.URL,
			Interval: target.Interval,
			Timeout:  target.Timeout,
		})
	}

	go prober.RunProbes(ctx)

	srv := server.NewServer(prober, exporter)

	log.Printf("Starting server on :%d", cfg.Server.Port)
	if err := srv.Start(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
