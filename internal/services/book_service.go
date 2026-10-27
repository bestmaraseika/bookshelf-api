package services

import (
	"context"

	"bookshelf/internal/models"
)

type BookRepository interface {
	Create(ctx context.Context, book *models.Book) error
	GetByID(ctx context.Context, id uint64) (*models.Book, error)
	GetAll(ctx context.Context, page, limit int) ([]models.Book, error)
	Update(ctx context.Context, book *models.Book) error
	Delete(ctx context.Context, id uint64) error
}

type BookService struct {
	repo BookRepository
}

func NewBookService(repo BookRepository) *BookService {
	return &BookService{
		repo: repo,
	}
}

func (s *BookService) Create(
	ctx context.Context,
	book *models.Book,
) error {
	return s.repo.Create(ctx, book)
}

func (s *BookService) GetByID(
	ctx context.Context,
	id uint64,
) (*models.Book, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *BookService) GetAll(
	ctx context.Context, page, limit int,
) ([]models.Book, error) {
	return s.repo.GetAll(ctx, page, limit)
}

func (s *BookService) Update(
	ctx context.Context, book *models.Book,
) error {
	return s.repo.Update(ctx, book)
}

func (s *BookService) Delete(
	ctx context.Context, id uint64,
) error {
	return s.repo.Delete(ctx, id)
}
