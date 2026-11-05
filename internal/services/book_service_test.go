package services

import (
	"bookshelf/internal/models"
	"context"
	"errors"
	"testing"
)

type mockBookRepository struct {
	getRecommendedFunc func(
		ctx context.Context,
		limit int,
	) ([]models.Book, error)

	markOutOfStockFunc func(
		ctx context.Context,
		id uint64,
	) error
}

func (m *mockBookRepository) GetRecommended(
	ctx context.Context,
	limit int,
) ([]models.Book, error) {
	return m.getRecommendedFunc(ctx, limit)
}

func (m *mockBookRepository) MarkOutOfStock(
	ctx context.Context,
	id uint64,
) error {
	return m.markOutOfStockFunc(ctx, id)
}

func (m *mockBookRepository) Create(
	ctx context.Context,
	book *models.Book,
) error {
	panic("unexpected Create call")
}

func (m *mockBookRepository) GetByID(
	ctx context.Context,
	id uint64,
) (*models.Book, error) {
	panic("unexpected GetByID call")
}

func (m *mockBookRepository) GetAll(
	ctx context.Context,
	page, limit int,
) ([]models.Book, error) {
	panic("unexpected GetAll call")
}

func (m *mockBookRepository) Update(
	ctx context.Context,
	book *models.Book,
) error {
	panic("unexpected Update call")
}

func (m *mockBookRepository) Delete(
	ctx context.Context,
	id uint64,
) error {
	panic("unexpected Delete call")
}

func TestBookService_Recommend(t *testing.T) {
	called := false

	mockRepo := &mockBookRepository{
		getRecommendedFunc: func(
			ctx context.Context,
			limit int,
		) ([]models.Book, error) {
			called = true

			if limit != 5 {
				t.Errorf("expected limit 5, got %d", limit)
			}

			return []models.Book{
				{ID: 1, Title: "Dune"},
				{ID: 2, Title: "Foundation"},
			}, nil
		},
	}

	service := NewBookService(mockRepo)

	books, err := service.Recommend(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !called {
		t.Fatal("GetRecommended was not called")
	}

	if len(books) != 2 {
		t.Fatalf("expected 2 books, got %d", len(books))
	}

	if books[0].Title != "Dune" {
		t.Errorf("unexpected first book: %s", books[0].Title)
	}
}

func TestBookService_MarkOutOfStock(t *testing.T) {
	var receivedID uint64
	called := false

	mockRepo := &mockBookRepository{
		markOutOfStockFunc: func(
			ctx context.Context,
			id uint64,
		) error {
			called = true
			receivedID = id
			return nil
		},
	}

	service := NewBookService(mockRepo)

	err := service.MarkOutOfStock(
		context.Background(),
		42,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !called {
		t.Fatal("MarkOutOfStock was not called")
	}

	if receivedID != 42 {
		t.Errorf("expected ID 42, got %d", receivedID)
	}
}

func TestBookService_MarkOutOfStock_Error(t *testing.T) {
	expectedErr := errors.New("database unavailable")

	mockRepo := &mockBookRepository{
		markOutOfStockFunc: func(
			ctx context.Context,
			id uint64,
		) error {
			return expectedErr
		},
	}

	service := NewBookService(mockRepo)

	err := service.MarkOutOfStock(
		context.Background(),
		42,
	)

	if !errors.Is(err, expectedErr) {
		t.Errorf(
			"expected error %v, got %v",
			expectedErr,
			err,
		)
	}
}
