package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/c-j-p-nordquist/ekolod/pkg/probe"
)

func main() {
	target := probe.Target{
		Name:     "Example",
		URL:      "https://example.com",
		Interval: 10 * time.Second,
		Timeout:  5 * time.Second,
	}

	ticker := time.NewTicker(target.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			result := probe.ProbeTarget(target)
			jsonResult, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(jsonResult))
		}
	}
}
