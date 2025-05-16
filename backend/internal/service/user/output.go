package user

import (
	"github.com/google/uuid"
	"time"
)

type UserOut struct {
	ID        int64      `json:"-"`
	UID       uuid.UUID  `json:"uid"`
	Email     string     `json:"email"`
	Name      string     `json:"name" `
	Role      string     `json:"role"`
	Password  string     `json:"password"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy string     `json:"created_by"`
	UpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy string     `json:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at"`
	DeletedBy string     `json:"deleted_by"`
}

func (u *UserOut) PopulateFromEntity(e *User) {
	u.ID = e.ID
	u.UID = e.UID
	u.Email = e.Email.String
	u.Name = e.Name.String
	u.Role = e.Role.String
	u.Password = e.Password.String
	u.CreatedAt = e.CreatedAt.Time
	u.CreatedBy = e.CreatedBy.String
	if !e.UpdatedAt.Time.IsZero() {
		u.UpdatedAt = &e.UpdatedAt.Time
	}
	u.UpdatedBy = e.UpdatedBy.String
	if !e.DeletedAt.Time.IsZero() {
		u.DeletedAt = &e.DeletedAt.Time
	}
	u.DeletedBy = e.DeletedBy.String
}
