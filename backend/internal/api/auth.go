package api

import (
	authmiddleware "backend/internal/service/middleware"
	auth "backend/internal/usecase/auth"
	"backend/pkg/chi/middleware"
	"backend/pkg/chi/response"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
)

// RegisterHandlers registers handlers for different HTTP requests.
func RegisterAuthHandlers(r chi.Router, usecase auth.Usecase, authMiddleware authmiddleware.AuthMiddleware, logger zerolog.Logger) {
	res := authResource{usecase, logger}
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", middleware.APIWrapper(res.login))

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.HandleToken())
			r.Get("/user-info", middleware.APIWrapper(res.userInfo))
		})
	})

}

type authResource struct {
	usecase auth.Usecase
	logger  zerolog.Logger
}

// login returns a handler that handles user login request.
func (r authResource) login(_ http.ResponseWriter, req *http.Request) response.Body {
	ctx := req.Context()
	var input auth.LoginIn
	if err := render.Bind(req, &input); err != nil {
		r.logger.Error().Err(err).Msg("")
		return response.InvalidInput(err)
	}

	resp, err := r.usecase.Login(ctx, input)
	if err != nil {
		r.logger.Error().Err(err).Msg("")
		return response.InternalServerError(err.Error())
	}

	return response.OK("success", resp)
}

func (r authResource) userInfo(w http.ResponseWriter, req *http.Request) response.Body {
	resp, err := r.usecase.UserInfo(req.Context())
	if err != nil {
		r.logger.Error().Err(err).Msg("")
		return response.InternalServerError(err.Error())
	}

	return response.OK("success", resp)
}
