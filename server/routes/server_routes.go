package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/moevm/nosql2h24-transcribtion/handlers"
)

func ServerRoutes(r chi.Router) {
	r.Get("/servers", handlers.GetServers)
	r.Get("/servers/{id}", handlers.GetServerByID)
	r.Post("/servers", handlers.CreateServer)
	r.Put("/servers/{id}", handlers.UpdateServer)
	r.Patch("/servers/{id}", handlers.PatchServer)
	r.Delete("/servers/{id}", handlers.DeleteServer)

	r.Get("/servers/{id}/currentJobs", handlers.GetServerCurrentJobs)
	r.Get("/servers/{id}/completedJobs", handlers.GetServerCompletedJobs)
	//В целом следующий роут не нужен, все сразу делается через юзера
	r.Post("/servers/{id}/jobs/{job_id}", handlers.AddJobToServer)

}
