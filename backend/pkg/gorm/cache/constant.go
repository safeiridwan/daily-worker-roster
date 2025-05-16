// Package cache provides caching mechanisms for storing and retrieving data.
package cache

// contextKey is a custom type used for context keys to avoid collisions.
type contextKey int

const (
	// TableNameKey is the key used to store and retrieve table cache from context.
	TableNameKey contextKey = iota

	// CacheKeyPrefix is the prefix used for cache keys related to tables.
	CacheKeyPrefix = "TABLE:%s-"

	// CacheKeyFormat is the format string used to create cache keys.
	// It combines the CacheKeyPrefix with an additional identifier.
	CacheKeyFormat = CacheKeyPrefix + "%s"
)
