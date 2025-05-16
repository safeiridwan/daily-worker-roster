package auth

import "backend/internal/service/user"

type LoginOut struct {
	Token string `json:"token"`
}

type InfoOut struct {
	UID   string `json:"uid"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

func (o *InfoOut) Populate(e *user.UserOut) {
	o.UID = e.UID.String()
	o.Email = e.Email
	o.Name = e.Name
	o.Role = e.Role
}
