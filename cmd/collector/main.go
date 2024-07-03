package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/c-j-p-nordquist/ekolod/internal/collector"
	"github.com/c-j-p-nordquist/ekolod/pkg/health"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	// Load configuration from environment variables
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	collectorPort := os.Getenv("COLLECTOR_PORT")
	if collectorPort == "" {
		collectorPort = "8081" // Default to 8081 if not set
	}

	// Retry connection to the database
	var db *pgxpool.Pool
	var err error
	for i := 0; i < 30; i++ { // Try for 5 minutes
		db, err = pgxpool.Connect(context.Background(), dbURL)
		if err == nil {
			// Try to create the table
			_, err = db.Exec(context.Background(), `
				CREATE TABLE IF NOT EXISTS metrics (
					time TIMESTAMPTZ NOT NULL,
					target TEXT NOT NULL,
					check_type TEXT NOT NULL,
					duration DOUBLE PRECISION NOT NULL,
					success BOOLEAN NOT NULL,
					message TEXT,
					status_code INTEGER,
					content_length BIGINT,
					tls_version TEXT,
					cert_expiry_days INTEGER
				);
				SELECT create_hypertable('metrics', 'time', if_not_exists => TRUE);
			`)
			if err == nil {
				break
			}
		}
		log.Printf("Failed to connect to database or create table. Retrying in 10 seconds...")
		time.Sleep(10 * time.Second)
	}
	if err != nil {
		log.Fatalf("Unable to connect to database or create table after multiple attempts: %v", err)
	}
	defer db.Close()

	log.Println("Successfully connected to the database and ensured table exists")

	// Initialize health checker
	healthChecker := health.New()
	healthChecker.AddChecker(collector.NewDatabaseChecker(db))

	// Set up HTTP routes
	http.HandleFunc("/health", healthChecker.Handler())
	http.HandleFunc("/metrics", collector.MetricsHandler(db))
	http.HandleFunc("/timeseries", collector.TimeSeriesHandler(db))

	// Start the server
	log.Printf("Starting Collector server on :%s", collectorPort)
	log.Fatal(http.ListenAndServe(":"+collectorPort, nil))
}
