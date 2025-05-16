// Package cache provides caching mechanisms for storing and retrieving data.
package cache

import (
	"context"

	"github.com/go-gorm/caches/v4"
)

// CacheRepository defines the interface for cache operations
type CacheRepository interface {
	// Get retrieves a cached query for the given key
	Get(ctx context.Context, key string, q *caches.Query[any]) (*caches.Query[any], error)

	// Store caches a query for the given key
	Store(ctx context.Context, key string, val *caches.Query[any]) error

	// Invalidate clears the cache
	Invalidate(ctx context.Context) error
}

// KeyGenerator defines the interface for generating cache keys
type KeyGenerator interface {
	// GenerateKey creates a cache key based on the given context and key string
	GenerateKey(ctx context.Context, key string) (string, error)
}
