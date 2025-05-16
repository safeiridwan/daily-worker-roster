package auth

import (
	"backend/pkg/chi/request"
	"net/http"
)

type LoginIn struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (r LoginIn) Bind(_ *http.Request) error {
	return request.ValidateStruct(r)
}
