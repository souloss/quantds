package manager

import (
	"context"
	"time"
)

// Storage defines the interface for optional data persistence.
// The Manager can use a Storage implementation as a persistent layer
// below the in-memory cache, allowing data to survive process restarts.
type Storage interface {
	// Save stores data with the given key and TTL.
	Save(ctx context.Context, key string, data []byte, ttl time.Duration) error
	// Load retrieves data by key. Returns an error if the key does not exist or has expired.
	Load(ctx context.Context, key string) ([]byte, error)
	// Delete removes data by key.
	Delete(ctx context.Context, key string) error
	// Close releases any resources held by the storage.
	Close() error
}
