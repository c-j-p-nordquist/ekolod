-- Enable the TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Create the metrics table if it doesn't exist
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

-- Create the hypertable
SELECT create_hypertable('metrics', 'time', if_not_exists => TRUE);