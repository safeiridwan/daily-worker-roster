// Package cache provides caching mechanisms for storing and retrieving data.
package cache

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// DefaultKeyGenerator is the default implementation of KeyGenerator.
type DefaultKeyGenerator struct{}

// GenerateKey generates a unique cache key based on the provided context and key string.
// It uses the table name from the context and an MD5 hash of the key to create a unique identifier.
//
// Parameters:
//   - ctx: The context containing the table name.
//   - key: The original key string to be hashed.
//
// Returns:
//   - string: The generated cache key.
//   - error: An error if the table name cannot be retrieved from the context.
func (g *DefaultKeyGenerator) GenerateKey(ctx context.Context, key string) (string, error) {
	table, ok := ctx.Value(TableNameKey).(string)
	if !ok {
		return "", nil
	}
	hash := md5.Sum([]byte(key))
	return fmt.Sprintf(CacheKeyFormat, table, hex.EncodeToString(hash[:])), nil
}
