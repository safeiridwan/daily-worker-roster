// Package dbcontext provides a wrapper around gorm.DB with additional functionality
// for handling database transactions and context management.
package dbcontext

import (
	"backend/pkg/gorm/cache"
	"context"
	"net/http"

	"gorm.io/gorm"
)

// sessions represents a database context that wraps a gorm.DB instance.
type sessions struct {
	db *gorm.DB
}

// New creates a new Sessions instance with the provided gorm.DB.
//
// Parameters:
//   - db: A pointer to a gorm.DB instance.
//
// Returns:
//   - Sessions: A new Sessions instance.
func New(db *gorm.DB) Sessions {
	return &sessions{db: db}
}

// With returns a new gorm.DB instance with the provided context and optional table cache.
// If a transaction is found in the context, it will be used instead of the default database connection.
//
// Parameters:
//   - ctx: The context.Context to be used.
//   - tableNameCache: Optional string slice for table caching.
//
// Returns:
//   - *gorm.DB: A new gorm.DB instance with the provided context.
func (d *sessions) With(ctx context.Context, tableNameCache ...string) *gorm.DB {
	var dbConn *gorm.DB
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {
		dbConn = tx
	} else {
		dbConn = d.db
	}

	if len(tableNameCache) > 0 {
		return dbConn.WithContext(context.WithValue(ctx, cache.TableNameKey, tableNameCache[0]))
	}
	return dbConn.WithContext(ctx)
}

// Transactional executes the provided function within a database transaction.
// If the function returns an error, the transaction is rolled back. Otherwise, it is committed.
//
// Parameters:
//   - ctx: The context.Context to be used.
//   - f: A function that takes a context.Context and returns an error.
//
// Returns:
//   - error: Any error that occurred during the transaction.
func (d *sessions) Transactional(ctx context.Context, f func(ctx context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return f(context.WithValue(ctx, txKey, tx))
	})
}

// TransactionHandler returns a middleware function that wraps the next http.Handler in a database transaction.
// If an error occurs during the transaction, it responds with an Internal Server Error.
//
// Returns:
//   - func(next http.Handler) http.Handler: A middleware function for handling database transactions.
func (d *sessions) TransactionHandler() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := d.db.WithContext(r.Context()).Transaction(func(tx *gorm.DB) error {
				ctx := context.WithValue(r.Context(), txKey, tx)
				next.ServeHTTP(w, r.WithContext(ctx))
				return nil
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
	}
}
