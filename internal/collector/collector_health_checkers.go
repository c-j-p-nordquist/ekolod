package collector

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DatabaseChecker struct {
	db *pgxpool.Pool
}

func NewDatabaseChecker(db *pgxpool.Pool) *DatabaseChecker {
	return &DatabaseChecker{db: db}
}

func (c *DatabaseChecker) Check() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result int
	err := c.db.QueryRow(ctx, "SELECT 1").Scan(&result)
	return err == nil && result == 1
}

func (c *DatabaseChecker) Name() string {
	return "database_connection"
}
