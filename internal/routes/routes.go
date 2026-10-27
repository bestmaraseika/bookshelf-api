package routes

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"

	"bookshelf/internal/handlers"
	"bookshelf/internal/repository"
	"bookshelf/internal/services"
)

func NewRouter(db *sql.DB) http.Handler {
	router := chi.NewRouter()

	bookRepo := repository.NewBookRepository(db)
	bookService := services.NewBookService(bookRepo)
	bookHandler := handlers.NewBookHandler(bookService)

	router.Get("/health", handlers.Health)

	router.Post("/books", bookHandler.CreateBook)
	router.Get("/books", bookHandler.GetBooks)
	router.Get("/books/{id}", bookHandler.GetBookByID)
	router.Put("/books/{id}", bookHandler.UpdateBook)
	router.Delete("/books/{id}", bookHandler.DeleteBook)

	return router
}
