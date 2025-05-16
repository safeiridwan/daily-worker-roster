// Package dbcontext provides a wrapper around gorm.DB with additional functionality
// for handling database transactions and context management.
package dbcontext

import (
	"context"
	"net/http"

	"gorm.io/gorm"
)

// Sessions represents an interface for database context operations
type Sessions interface {
	// With returns a new *gorm.DB instance with the given context and optional table cache
	With(ctx context.Context, tableCache ...string) *gorm.DB

	// Transactional executes the given function within a database transaction
	Transactional(ctx context.Context, f func(ctx context.Context) error) error

	// TransactionHandler returns an http.Handler middleware for managing database transactions
	TransactionHandler() func(next http.Handler) http.Handler
}
