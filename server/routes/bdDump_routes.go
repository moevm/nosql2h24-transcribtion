package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/moevm/nosql2h24-transcribtion/handlers"
)

func bdDumpRoutes(r chi.Router) {
	r.Get("/dump/export", handlers.ExportData)
	r.Post("/dump/import", handlers.ImportData)

}
