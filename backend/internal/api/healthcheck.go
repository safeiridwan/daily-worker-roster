package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

func RegisterHandlers(r *chi.Mux, version string) {
	r.HandleFunc("GET /healthcheck", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, "OK "+version)
	})
}
