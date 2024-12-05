package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/moevm/nosql2h24-transcribtion/handlers"
)

func ServerRoutes(r chi.Router) {
	r.Get("/servers", handlers.GetServers) // Получение списка серверов
	//r.Get("/servers/{id}", handlers.GetServerByID)   // Получение сервера по ID
	//r.Post("/servers", handlers.CreateServer)        // Создание нового сервера
	//r.Put("/servers/{id}", handlers.UpdateServer)    // Полное обновление данных сервера
	//r.Patch("/servers/{id}", handlers.PatchServer)   // Частичное обновление данных сервера
	//r.Delete("/servers/{id}", handlers.DeleteServer) // Удаление сервера
	//
	//r.Get("/servers/{id}/tasks", handlers.GetServerTasks)                   // Список текущих заданий сервера
	//r.Post("/servers/{id}/tasks", handlers.AddServerTask)                   // Добавление задания серверу
	//r.Delete("/servers/{id}/tasks/{task_id}", handlers.RemoveServerTask)    // Удаление задания с сервера
	//r.Post("/servers/{id}/tasks/{task_id}/complete", handlers.CompleteTask) // Завершение задания
	//
	//r.Patch("/servers/{id}/status", handlers.UpdateServerStatus) // Обновление статуса сервера
}
