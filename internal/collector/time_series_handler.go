package collector

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func TimeSeriesHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		target := r.URL.Query().Get("target")
		checkType := r.URL.Query().Get("check_type")
		duration := r.URL.Query().Get("duration")

		// Set default duration if not provided
		if duration == "" {
			duration = "1h" // Default to last 1 hour
		}

		// Parse duration
		dur, err := time.ParseDuration(duration)
		if err != nil {
			http.Error(w, "Invalid duration", http.StatusBadRequest)
			return
		}

		// Calculate time range
		endTime := time.Now()
		startTime := endTime.Add(-dur)

		// Construct and execute the query
		query := `
			SELECT time, target, check_type, duration, success
			FROM metrics
			WHERE time BETWEEN $1 AND $2
			AND ($3 = '' OR target = $3)
			AND ($4 = '' OR check_type = $4)
			ORDER BY time ASC
		`
		rows, err := db.Query(r.Context(), query, startTime, endTime, target, checkType)
		if err != nil {
			http.Error(w, "Failed to fetch time series data", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Process the results
		var results []map[string]interface{}
		for rows.Next() {
			var (
				timestamp time.Time
				target    string
				checkType string
				duration  float64
				success   bool
			)
			if err := rows.Scan(&timestamp, &target, &checkType, &duration, &success); err != nil {
				http.Error(w, "Failed to process time series data", http.StatusInternalServerError)
				return
			}
			results = append(results, map[string]interface{}{
				"time":       timestamp,
				"target":     target,
				"check_type": checkType,
				"duration":   duration,
				"success":    success,
			})
		}

		// Send the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}
