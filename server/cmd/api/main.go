package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/moevm/nosql2h24-transcribtion/config"
	"github.com/moevm/nosql2h24-transcribtion/db"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	client := db.InitConnection(&cfg)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := client.Ping(ctx, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error pinging database", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	if err := http.ListenAndServe(cfg.Port, r); err != nil {
		log.Fatal("Could not start http server ", err)
	}

	fmt.Println("Went past server")
}
