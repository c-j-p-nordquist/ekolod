package collector

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func MetricsHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payload struct {
			Target string `json:"target"`
			Check  string `json:"check"`
			Result struct {
				Duration       float64 `json:"duration"`
				Success        bool    `json:"success"`
				Message        string  `json:"message"`
				StatusCode     int     `json:"statusCode"`
				ContentLength  int64   `json:"contentLength"`
				TLSVersion     string  `json:"tlsVersion"`
				CertExpiryDays int     `json:"certExpiryDays"`
			} `json:"result"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		_, err := db.Exec(r.Context(),
			`INSERT INTO metrics (time, target, check_type, duration, success, message, status_code, content_length, tls_version, cert_expiry_days)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			time.Now(), payload.Target, payload.Check,
			payload.Result.Duration, payload.Result.Success, payload.Result.Message,
			payload.Result.StatusCode, payload.Result.ContentLength,
			payload.Result.TLSVersion, payload.Result.CertExpiryDays)

		if err != nil {
			http.Error(w, "Failed to insert metrics", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
