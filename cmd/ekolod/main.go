package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/c-j-p-nordquist/ekolod/pkg/probe"
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

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			results := prober.GetResults()
			jsonResult, _ := json.MarshalIndent(results, "", "  ")
			fmt.Println(string(jsonResult))
		}
	}
}
