package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() chi.Router {

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	UserRoutes(r)

	return r
}
