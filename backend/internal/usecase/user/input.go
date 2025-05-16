package user

import (
	"backend/pkg/chi/request"
	"net/http"
)

type UserIn struct {
	Email string `json:"email"`
	Name  string `json:"name" `
	Role  string `json:"role" `
}

func (c UserIn) Bind(_ *http.Request) error {
	return request.ValidateStruct(c)
}
