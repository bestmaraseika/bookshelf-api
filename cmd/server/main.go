package main

import (
	"log"
	"net/http"

	"bookshelf/internal/config"
	"bookshelf/internal/repository"
	"bookshelf/internal/routes"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := repository.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("Connected to PostgreSQL")

	router := routes.NewRouter(db)

	addr := ":" + cfg.AppPort

	log.Println("BookShelf API started on", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
