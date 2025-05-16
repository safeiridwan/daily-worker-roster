package api

import (
	authmiddleware "backend/internal/service/middleware"
	scheduleusecase "backend/internal/usecase/schedule"
	"backend/pkg/chi/middleware"
	"backend/pkg/chi/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"net/http"
)

func RegisterScheduleHandlers(r chi.Router, usecase scheduleusecase.Usecase, authmiddleware authmiddleware.AuthMiddleware, logger zerolog.Logger) {
	res := scheduleResource{usecase: usecase, logger: logger}
	r.Route("/schedule", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(authmiddleware.HandleToken())
			r.Post("/", middleware.APIWrapper(res.createSchedule))
		})
	})
}

type scheduleResource struct {
	usecase scheduleusecase.Usecase
	logger  zerolog.Logger
}

func (r scheduleResource) createSchedule(_ http.ResponseWriter, req *http.Request) response.Body {
	ctx := req.Context()
	identity := authmiddleware.CurrentUser(ctx)
	identity.UserUID = "uuid"
	var input scheduleusecase.ScheduleIn
	if err := render.Bind(req, &input); err != nil {
		r.logger.Error().Err(err).Msg("")
		return response.InvalidInput(err)
	}

	if err := r.usecase.CreateSchedule(ctx, input, *identity); err != nil {
		r.logger.Error().Err(err).Msg("")
		return response.InternalServerError(err.Error())
	}

	return response.OK("Successfully created schedule", nil)
}
