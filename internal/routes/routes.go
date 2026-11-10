package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"bookshelf/internal/handlers"
	"bookshelf/internal/repository"
	"bookshelf/internal/services"
)

func NewRouter(db *sql.DB) http.Handler {
	router := chi.NewRouter()

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		writeRoutingError(w, http.StatusNotFound, "route not found")
	})
	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		for _, method := range []string{
			http.MethodConnect, http.MethodDelete, http.MethodGet, http.MethodHead,
			http.MethodOptions, http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodTrace,
		} {
			if router.Match(chi.NewRouteContext(), method, r.URL.Path) {
				w.Header().Add("Allow", method)
			}
		}
		writeRoutingError(w, http.StatusMethodNotAllowed, "method not allowed")
	})

	bookRepo := repository.NewBookRepository(db)
	bookService := services.NewBookService(bookRepo)
	bookHandler := handlers.NewBookHandler(bookService)

	router.Get("/health", handlers.Health)

	router.Post("/books", bookHandler.CreateBook)
	router.Get("/books", bookHandler.GetBooks)
	router.Get("/books/recommend", bookHandler.RecommendBooks)
	router.Get("/books/{id}", bookHandler.GetBookByID)
	router.Put("/books/{id}", bookHandler.UpdateBook)
	router.Post("/books/{id}/mark-out-of-stock", bookHandler.MarkOutOfStock)
	router.Delete("/books/{id}", bookHandler.DeleteBook)

	return router
}

func writeRoutingError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
