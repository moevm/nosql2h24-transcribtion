package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/moevm/nosql2h24-transcribtion/handlers"
)

func UserRoutes(r chi.Router) {
	r.Get("/users", handlers.GetUsers)
}
