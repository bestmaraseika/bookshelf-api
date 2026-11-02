package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"bookshelf/internal/models"
	"bookshelf/internal/services"
)

type BookHandler struct {
	service *services.BookService
}

func NewBookHandler(service *services.BookService) *BookHandler {
	return &BookHandler{
		service: service,
	}
}

type createBookRequest struct {
	Title      string  `json:"title"`
	Author     string  `json:"author"`
	Year       int     `json:"year"`
	ISBN       *string `json:"isbn"`
	OutOfStock bool    `json:"out_of_stock"`
	Rating     *int    `json:"rating"`
}

func (h *BookHandler) CreateBook(
	w http.ResponseWriter,
	r *http.Request,
) {
	req, err := decodeBookRequest(w, r)
	if err != nil {
		handleDecodeError(w, err)
		return
	}

	if err := validateBookRequest(req); err != nil {
		writeJSONError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	book := &models.Book{
		Title:      req.Title,
		Author:     req.Author,
		Year:       req.Year,
		ISBN:       req.ISBN,
		OutOfStock: req.OutOfStock,
		Rating:     req.Rating,
	}

	if err := h.service.Create(r.Context(), book); err != nil {
		handleDatabaseError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, book)
}

func (h *BookHandler) GetBookByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	rawID := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 || id > 9223372036854775807 {
		writeJSONError(w, http.StatusBadRequest, "invalid book ID")
		return
	}

	book, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSONError(w, http.StatusNotFound, "book not found")
			return
		}

		handleDatabaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 10

	if value := r.URL.Query().Get("page"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid page")
			return
		}
		page = parsed
	}

	if value := r.URL.Query().Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid limit")
			return
		}
		limit = parsed
	}

	if page < 1 || limit < 1 || limit > 100 {
		writeJSONError(w, http.StatusBadRequest, "invalid pagination")
		return
	}

	books, err := h.service.GetAll(r.Context(), page, limit)
	if err != nil {
		handleDatabaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, books)

}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id == 0 || id > 9223372036854775807 {
		writeJSONError(w, http.StatusBadRequest, "invalid book ID")
		return
	}

	req, err := decodeBookRequest(w, r)
	if err != nil {
		handleDecodeError(w, err)
		return
	}

	if err := validateBookRequest(req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	book := &models.Book{
		ID:         id,
		Title:      req.Title,
		Author:     req.Author,
		Year:       req.Year,
		ISBN:       req.ISBN,
		OutOfStock: req.OutOfStock,
		Rating:     req.Rating,
	}

	if err := h.service.Update(r.Context(), book); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSONError(w, http.StatusNotFound, "book not found")
			return
		}
		handleDatabaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id == 0 || id > 9223372036854775807 {
		writeJSONError(w, http.StatusBadRequest, "invalid book ID")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSONError(w, http.StatusNotFound, "book not found")
			return
		}
		handleDatabaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *BookHandler) RecommendBooks(
	w http.ResponseWriter,
	r *http.Request,
) {
	books, err := h.service.Recommend(r.Context())
	if err != nil {
		handleDatabaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, books)
}

func (h *BookHandler) MarkOutOfStock(
	w http.ResponseWriter,
	r *http.Request,
) {
	rawID := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 || id > 9223372036854775807 {
		writeJSONError(w, http.StatusBadRequest, "invalid book ID")
		return
	}

	err = h.service.MarkOutOfStock(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSONError(w, http.StatusNotFound, "book not found")
			return
		}

		handleDatabaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}
