package postgres

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"

	"src/backend/internal/config"
)

// Store wraps a pgxpool.Pool and implements repository.Store.
// All repository methods are defined in their own files within this package.
type Store struct {
	pool *pgxpool.Pool
}

// New creates a connection pool, verifies connectivity with an initial Ping,
// and returns a ready-to-use Store.
// ctx controls the timeout for the initial connection attempt.
func New(ctx context.Context, cfg config.DatabaseConfig) (*Store, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		url.QueryEscape(cfg.User),
		url.QueryEscape(cfg.Password),
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &Store{pool: pool}, nil
}

// Close releases all pool connections. Call this on application shutdown.
func (s *Store) Close() {
	s.pool.Close()
}
