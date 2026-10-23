package repository

import (
	"bookshelf/internal/models"
	"context"
	"database/sql"
	"fmt"
	"math"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{
		db: db,
	}
}

func (r *BookRepository) Create(
	ctx context.Context,
	book *models.Book,
) error {
	query := `
        INSERT INTO books (
            title,
            author,
            year,
            isbn,
            out_of_stock,
            rating
        )
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at
    `

	err := r.db.QueryRowContext(
		ctx,
		query,
		book.Title,
		book.Author,
		book.Year,
		book.ISBN,
		book.OutOfStock,
		book.Rating,
	).Scan(
		&book.ID,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("create book: %w", err)
	}

	return nil
}

func (r *BookRepository) GetByID(
	ctx context.Context,
	id uint64,
) (*models.Book, error) {
	query := `
        SELECT
            id,
            title,
            author,
            year,
            isbn,
            out_of_stock,
            rating,
            created_at,
            updated_at
        FROM books
        WHERE id = $1
    `

	if id > 9223372036854775807 {
		return nil, fmt.Errorf("book ID is out of range")
	}

	var book models.Book

	err := r.db.QueryRowContext(
		ctx,
		query,
		int64(id),
	).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Year,
		&book.ISBN,
		&book.OutOfStock,
		&book.Rating,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("get book by ID: %w", err)
	}

	return &book, nil
}

func (r *BookRepository) GetAll(
	ctx context.Context,
	page int,
	limit int,
) ([]models.Book, error) {
	if page < 1 || limit < 1 || limit > 100 {
		return nil, fmt.Errorf("invalid pagination parameters")
	}

	if int64(page-1) > (math.MaxInt64-int64(limit))/int64(limit) {
		return nil, fmt.Errorf("pagination offset is too large")
	}

	offset := int64(page-1) * int64(limit)

	query := `
        SELECT
            id,
            title,
            author,
            year,
            isbn,
            out_of_stock,
            rating,
            created_at,
            updated_at
        FROM books
        ORDER BY id ASC
        LIMIT $1 OFFSET $2
    `

	rows, err := r.db.QueryContext(
		ctx,
		query,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("get books: %w", err)
	}
	defer rows.Close()

	books := make([]models.Book, 0)

	for rows.Next() {
		var book models.Book

		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Year,
			&book.ISBN,
			&book.OutOfStock,
			&book.Rating,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan book: %w", err)
		}

		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate books: %w", err)
	}

	return books, nil
}

func (r *BookRepository) Update(
	ctx context.Context,
	book *models.Book,
) error {
	if book.ID > 9223372036854775807 {
		return fmt.Errorf("book ID is out of range")
	}

	query := `
        UPDATE books
        SET
            title = $1,
            author = $2,
            year = $3,
            isbn = $4,
            out_of_stock = $5,
            rating = $6,
            updated_at = NOW()
        WHERE id = $7
        RETURNING updated_at
    `

	err := r.db.QueryRowContext(
		ctx,
		query,
		book.Title,
		book.Author,
		book.Year,
		book.ISBN,
		book.OutOfStock,
		book.Rating,
		int64(book.ID),
	).Scan(&book.UpdatedAt)

	if err != nil {
		return fmt.Errorf("update book: %w", err)
	}

	return nil
}

func (r *BookRepository) Delete(
	ctx context.Context,
	id uint64,
) error {
	if id > 9223372036854775807 {
		return fmt.Errorf("book ID is out of range")
	}

	query := `
        DELETE FROM books
        WHERE id = $1
    `

	result, err := r.db.ExecContext(
		ctx,
		query,
		int64(id),
	)
	if err != nil {
		return fmt.Errorf("delete book: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("delete book: %w", sql.ErrNoRows)
	}

	return nil
}
