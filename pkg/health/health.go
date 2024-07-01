package health

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type HealthChecker interface {
	Check() bool
	Name() string
}

type Health struct {
	checkers []HealthChecker
	mu       sync.RWMutex
}

type HealthStatus struct {
	Status    string                 `json:"status"`
	Checks    map[string]bool        `json:"checks"`
	Timestamp string                 `json:"timestamp"`
	Info      map[string]interface{} `json:"info,omitempty"`
}

func New() *Health {
	return &Health{
		checkers: []HealthChecker{},
	}
}

func (h *Health) AddChecker(checker HealthChecker) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checkers = append(h.checkers, checker)
}

func (h *Health) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		status := HealthStatus{
			Status:    "healthy",
			Checks:    make(map[string]bool),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}

		h.mu.RLock()
		for _, checker := range h.checkers {
			checkStatus := checker.Check()
			status.Checks[checker.Name()] = checkStatus
			if !checkStatus {
				status.Status = "unhealthy"
			}
		}
		h.mu.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		if status.Status == "unhealthy" {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		json.NewEncoder(w).Encode(status)
	}
}
