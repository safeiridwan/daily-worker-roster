package api

import (
	"backend/internal/usecase/user"
	"backend/pkg/chi/middleware"
	"backend/pkg/chi/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"net/http"
)

func RegisterUserHandlers(r chi.Router, usecase user.Usecase, logger zerolog.Logger) {
	res := userResource{usecase: usecase, logger: logger}
	r.Route("/user", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post("/", middleware.APIWrapper(res.createUser))
		})
	})
}

type userResource struct {
	usecase user.Usecase
	logger  zerolog.Logger
}

func (r userResource) createUser(_ http.ResponseWriter, req *http.Request) response.Body {
	ctx := req.Context()
	var input user.UserIn
	if err := render.Bind(req, &input); err != nil {
		r.logger.Error().Err(err).Msg("")
		return response.InvalidInput(err)
	}

	resp, err := r.usecase.CreateUser(ctx, input)
	if err != nil {
		r.logger.Error().Err(err).Msg("")
		return response.InternalServerError(err.Error())
	}

	return response.OK("Successfully created user", resp)
}
