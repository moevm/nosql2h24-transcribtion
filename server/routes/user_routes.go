package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/moevm/nosql2h24-transcribtion/handlers"
)

func UserRoutes(r chi.Router) {
	r.Get("/users", handlers.GetUsers)
	r.Get("/users/{id}", handlers.GetUserByID)

	r.Post("/users", handlers.CreateUser)
	r.Put("/users/{id}", handlers.UpdateUser)
	r.Patch("/users/{id}", handlers.PatchUser)

	r.Delete("/users/{id}", handlers.DeleteUser)

	r.Get("/users/{id}/jobs", handlers.GetUserJobs)
	r.Post("/users/{id}/jobs", handlers.AddUserJob)
	r.Delete("/users/{id}/jobs/{jobId}", handlers.DeleteUserJob)
}
