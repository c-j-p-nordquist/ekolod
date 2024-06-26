package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/c-j-p-nordquist/ekolod/pkg/config"
	"github.com/c-j-p-nordquist/ekolod/pkg/probe"
	"github.com/c-j-p-nordquist/ekolod/pkg/server"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
)

func main() {
	log.Println("Starting Ekolod probe application")

	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("Initial configuration loaded successfully")

	exporter, err := prometheus.New()
	if err != nil {
		log.Fatalf("Failed to create Prometheus exporter: %v", err)
	}
	log.Println("Prometheus exporter created successfully")

	meterProvider := sdk.NewMeterProvider(sdk.WithReader(exporter))
	meter := meterProvider.Meter("ekolod")

	prober := createProber(cfg, meter)
	log.Printf("Prober created with %d targets", len(cfg.Targets))

	srv := server.NewServer(prober, exporter)

	go prober.RunProbes()

	go func() {
		for {
			select {
			case <-srv.WaitForReload():
				log.Println("Reload signal received, reloading configuration...")
				newCfg, err := config.Load("config.yaml")
				if err != nil {
					log.Printf("Failed to reload config: %v", err)
					continue
				}
				log.Println("New configuration loaded successfully")

				newProber := createProber(newCfg, meter)
				log.Printf("New prober created with %d targets", len(newCfg.Targets))

				srv.UpdateProber(newProber)
				go newProber.RunProbes()

				log.Println("Configuration reloaded and applied successfully")
			}
		}
	}()

	// Set up graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("Received shutdown signal. Initiating graceful shutdown...")
		prober.Stop()
		os.Exit(0)
	}()

	log.Printf("Starting server on :%d", cfg.Server.Port)
	if err := srv.Start(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func createProber(cfg *config.Config, meter metric.Meter) *probe.Prober {
	prober := probe.NewProber(meter)
	for _, target := range cfg.Targets {
		prober.AddTarget(probe.Target{
			Name:     target.Name,
			URL:      target.URL,
			Interval: target.Interval,
			Timeout:  target.Timeout,
		})
		log.Printf("Added target: %s (%s)", target.Name, target.URL)
	}
	return prober
}
