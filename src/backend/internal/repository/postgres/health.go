package postgres

import "context"

// Ping verifies that the database is reachable by acquiring a connection
// from the pool and running a trivial server round-trip.
func (s *Store) Ping(ctx context.Context) error {
	return s.pool.Ping(ctx)
}
