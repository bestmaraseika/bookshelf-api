package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"bookshelf/internal/handlers"
)

func NewRouter() http.Handler {
	router := chi.NewRouter()

	router.Get("/health", handlers.Health)

	return router
}
