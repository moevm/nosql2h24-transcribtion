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

	//r.Get("/servers/{id}/tasks", handlers.GetServerTasks)                   // Список текущих заданий сервера
	//r.Post("/servers/{id}/tasks", handlers.AddServerTask)                   // Добавление задания серверу
	//r.Delete("/servers/{id}/tasks/{task_id}", handlers.RemoveServerTask)    // Удаление задания с сервера
	//r.Post("/servers/{id}/tasks/{task_id}/complete", handlers.CompleteTask) // Завершение задания
	//
	//r.Patch("/servers/{id}/status", handlers.UpdateServerStatus) // Обновление статуса сервера
}
