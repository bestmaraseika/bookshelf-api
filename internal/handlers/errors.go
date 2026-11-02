package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

func handleDatabaseError(
	w http.ResponseWriter,
	err error,
) {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			writeJSONError(
				w,
				http.StatusConflict,
				"resource already exists",
			)
			return

		case "23514", "23502":
			writeJSONError(
				w,
				http.StatusBadRequest,
				"database constraint violation",
			)
			return
		}
	}

	log.Printf("database error: %v", err)

	writeJSONError(
		w,
		http.StatusInternalServerError,
		"internal server error",
	)
}

func handleDecodeError(
	w http.ResponseWriter,
	err error,
) {
	var maxBytesErr *http.MaxBytesError

	if errors.As(err, &maxBytesErr) {
		writeJSONError(
			w,
			http.StatusRequestEntityTooLarge,
			"request body is too large",
		)
		return
	}

	log.Printf("invalid request body: %v", err)

	writeJSONError(
		w,
		http.StatusBadRequest,
		"invalid JSON body",
	)
}
