package routes

import (
	"github.com/go-chi/chi/v5"
)

// NewRouter создает и возвращает новый роутер
func NewRouter() chi.Router {
	r := chi.NewRouter()

	UserRoutes(r) // Пользовательские маршруты

	return r
}
