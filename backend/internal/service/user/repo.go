package user

import (
	"backend/pkg/gorm/dbcontext"
	"context"
	"database/sql"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        int64          `gorm:"primaryKey"`
	UID       uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid()"`
	Email     sql.NullString `gorm:"column:email"`
	Name      sql.NullString `gorm:"column:name"`
	Role      sql.NullString `gorm:"column:role"`
	Password  sql.NullString `gorm:"column:password"`
	CreatedAt sql.NullTime   `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	CreatedBy sql.NullString `gorm:"column:created_by"`
	UpdatedAt sql.NullTime   `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	UpdatedBy sql.NullString `gorm:"column:updated_by"`
	DeletedAt sql.NullTime   `gorm:"column:deleted_at"`
	DeletedBy sql.NullString `gorm:"column:deleted_by"`
}

func (User) TableName() string { return "user" }

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUID(ctx context.Context, uid string) (*User, error)
	SaveUser(ctx context.Context, user User) error
}

type userRepository struct {
	db dbcontext.Sessions
}

func NewUserRepository(db dbcontext.Sessions) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.With(ctx, "user").
		First(&user, "email = ?", email).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) GetByUID(ctx context.Context, uid string) (*User, error) {
	var user User
	err := r.db.With(ctx, "user").
		First(&user, "uid = ?", uid).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) SaveUser(ctx context.Context, user User) error {
	return r.db.With(ctx).Create(&user).Error
}
