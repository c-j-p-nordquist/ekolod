package main

import (
	"log"
	"time"

	"github.com/c-j-p-nordquist/ekolod/pkg/probe"
	"github.com/c-j-p-nordquist/ekolod/pkg/server"
)

func main() {
	prober := probe.NewProber()

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

	go prober.RunProbes()

	srv := server.NewServer(prober)
	log.Println("Starting server on :8080")
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
