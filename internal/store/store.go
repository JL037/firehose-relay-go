// Package store is the persistence side of the relay: it takes matched
// firehose events and writes them out in batches.
//
// SCAFFOLD: the interface seam is defined here so the pgx/Postgres writer can
// drop in during week 3 (see docs/planning). No implementation yet.
package store

import (
	"context"

	"github.com/JL037/firehose-relay-go/internal/firehose"
)

// Store persists a batch of matched events.
//
// TODO(week 3): implement a Postgres-backed store using
// github.com/jackc/pgx/v5, batching writes (benchmark COPY vs. INSERT).
type Store interface {
	Write(ctx context.Context, events []firehose.Event) error
	Close() error
}
