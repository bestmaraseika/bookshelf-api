package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"bookshelf/internal/models"
)

func TestBookRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create SQL mock: %v", err)
	}
	defer db.Close()

	repo := NewBookRepository(db)

	rating := 9
	now := time.Now().UTC()

	book := &models.Book{
		Title:  "Dune",
		Author: "Frank Herbert",
		Year:   1965,
		Rating: &rating,
	}

	rows := sqlmock.NewRows(
		[]string{"id", "created_at", "updated_at"},
	).AddRow(1, now, now)

	mock.ExpectQuery("INSERT INTO books").
		WithArgs(
			"Dune",
			"Frank Herbert",
			1965,
			nil,
			false,
			9,
		).
		WillReturnRows(rows)

	err = repo.Create(context.Background(), book)
	if err != nil {
		t.Fatalf("create book: %v", err)
	}

	if book.ID != 1 {
		t.Errorf("expected ID 1, got %d", book.ID)
	}

	if !book.CreatedAt.Equal(now) {
		t.Errorf("unexpected created_at: %v", book.CreatedAt)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet SQL expectations: %v", err)
	}
}

func TestBookRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create SQL mock: %v", err)
	}
	defer db.Close()

	repo := NewBookRepository(db)
	now := time.Now().UTC()

	rows := sqlmock.NewRows([]string{
		"id", "title", "author", "year",
		"isbn", "out_of_stock", "rating",
		"created_at", "updated_at",
	}).AddRow(
		7, "Dune", "Frank Herbert", 1965,
		nil, false, 9, now, now,
	)

	mock.ExpectQuery("SELECT").
		WithArgs(int64(7)).
		WillReturnRows(rows)

	book, err := repo.GetByID(context.Background(), 7)
	if err != nil {
		t.Fatalf("get book: %v", err)
	}

	if book.ID != 7 || book.Title != "Dune" {
		t.Errorf("unexpected book: %+v", book)
	}

	if book.Rating == nil || *book.Rating != 9 {
		t.Errorf("unexpected rating: %v", book.Rating)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet SQL expectations: %v", err)
	}
}

func TestBookRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create SQL mock: %v", err)
	}
	defer db.Close()

	repo := NewBookRepository(db)

	mock.ExpectExec("DELETE FROM books").
		WithArgs(int64(7)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(context.Background(), 7)
	if err != nil {
		t.Fatalf("delete book: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet SQL expectations: %v", err)
	}
}

func TestBookRepository_Delete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create SQL mock: %v", err)
	}
	defer db.Close()

	repo := NewBookRepository(db)

	mock.ExpectExec("DELETE FROM books").
		WithArgs(int64(999)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.Delete(context.Background(), 999)

	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet SQL expectations: %v", err)
	}
}
