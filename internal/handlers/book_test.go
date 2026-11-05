package handlers

import (
	"bookshelf/internal/models"
	"bookshelf/internal/services"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

type MockBookRepository struct {
	GetByIDFunc func(
		ctx context.Context,
		id uint64,
	) (*models.Book, error)
}

func (m *MockBookRepository) GetByID(
	ctx context.Context,
	id uint64,
) (*models.Book, error) {
	return m.GetByIDFunc(ctx, id)
}

func (m *MockBookRepository) Create(
	ctx context.Context,
	book *models.Book,
) error {
	panic("unexpected Create call")
}

func (m *MockBookRepository) GetAll(
	ctx context.Context,
	page, limit int,
) ([]models.Book, error) {
	panic("unexpected GetAll call")
}

func (m *MockBookRepository) Update(
	ctx context.Context,
	book *models.Book,
) error {
	panic("unexpected Update call")
}

func (m *MockBookRepository) Delete(
	ctx context.Context,
	id uint64,
) error {
	panic("unexpected Delete call")
}

func (m *MockBookRepository) GetRecommended(
	ctx context.Context,
	limit int,
) ([]models.Book, error) {
	panic("unexpected GetRecommended call")
}

func (m *MockBookRepository) MarkOutOfStock(
	ctx context.Context,
	id uint64,
) error {
	panic("unexpected MarkOutOfStock call")
}

func TestCreateBook_InvalidTitle(t *testing.T) {
	body := `{
        "title": "",
        "author": "Frank Herbert",
        "year": 1965,
        "rating": 9
    }`

	req := httptest.NewRequest(
		http.MethodPost,
		"/books",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	handler := NewBookHandler(nil)

	handler.CreateBook(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf(
			"expected status 400, got %d",
			recorder.Code,
		)
	}

	expectedBody := "{\"error\":\"title is required\"}\n"

	if recorder.Body.String() != expectedBody {
		t.Errorf(
			"unexpected response body: %s",
			recorder.Body.String(),
		)
	}
}

func TestGetBookByID_Success(t *testing.T) {
	expectedBook := &models.Book{
		ID:     1,
		Title:  "Dune",
		Author: "Frank Herbert",
		Year:   1965,
	}

	mockRepo := &MockBookRepository{
		GetByIDFunc: func(
			ctx context.Context,
			id uint64,
		) (*models.Book, error) {
			if id != 1 {
				t.Errorf("expected ID 1, got %d", id)
			}

			return expectedBook, nil
		},
	}

	service := services.NewBookService(mockRepo)
	handler := NewBookHandler(service)

	router := chi.NewRouter()
	router.Get("/books/{id}", handler.GetBookByID)

	req := httptest.NewRequest(
		http.MethodGet,
		"/books/1",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d: %s",
			recorder.Code,
			recorder.Body.String(),
		)
	}

	var actualBook models.Book

	err := json.Unmarshal(
		recorder.Body.Bytes(),
		&actualBook,
	)
	if err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if actualBook.ID != expectedBook.ID {
		t.Errorf(
			"expected ID %d, got %d",
			expectedBook.ID,
			actualBook.ID,
		)
	}

	if actualBook.Title != expectedBook.Title {
		t.Errorf(
			"expected title %q, got %q",
			expectedBook.Title,
			actualBook.Title,
		)
	}

	if actualBook.Author != expectedBook.Author {
		t.Errorf(
			"expected author %q, got %q",
			expectedBook.Author,
			actualBook.Author,
		)
	}
}

func TestGetBookByID_NotFound(t *testing.T) {
	mockRepo := &MockBookRepository{
		GetByIDFunc: func(
			ctx context.Context,
			id uint64,
		) (*models.Book, error) {
			return nil, sql.ErrNoRows
		},
	}

	service := services.NewBookService(mockRepo)
	handler := NewBookHandler(service)

	router := chi.NewRouter()
	router.Get("/books/{id}", handler.GetBookByID)

	req := httptest.NewRequest(
		http.MethodGet,
		"/books/999",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status 404, got %d",
			recorder.Code,
		)
	}

	expectedBody := "{\"error\":\"book not found\"}\n"

	if recorder.Body.String() != expectedBody {
		t.Errorf(
			"expected body %q, got %q",
			expectedBody,
			recorder.Body.String(),
		)
	}
}

func TestGetBookByID_InternalError(t *testing.T) {
	mockRepo := &MockBookRepository{
		GetByIDFunc: func(
			ctx context.Context,
			id uint64,
		) (*models.Book, error) {
			return nil, errors.New("database unavailable")
		},
	}

	service := services.NewBookService(mockRepo)
	handler := NewBookHandler(service)

	router := chi.NewRouter()
	router.Get("/books/{id}", handler.GetBookByID)

	req := httptest.NewRequest(
		http.MethodGet,
		"/books/1",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status 500, got %d",
			recorder.Code,
		)
	}

	expectedBody := "{\"error\":\"internal server error\"}\n"

	if recorder.Body.String() != expectedBody {
		t.Errorf(
			"expected body %q, got %q",
			expectedBody,
			recorder.Body.String(),
		)
	}
}
