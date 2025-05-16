// Package dbcontext provides a wrapper around gorm.DB with additional functionality
// for handling database transactions and context management.
package dbcontext

// contextKey is a custom type used for context keys to avoid collisions.
type contextKey int

const (
	// txKey is the context key for storing database transactions.
	txKey contextKey = iota
)
