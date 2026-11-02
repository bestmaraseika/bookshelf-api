package handlers

import (
	"errors"
	"strings"
)

func validateBookRequest(req createBookRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return errors.New("title is required")
	}

	if strings.TrimSpace(req.Author) == "" {
		return errors.New("author is required")
	}

	if req.Year <= 0 {
		return errors.New("year must be positive")
	}

	if req.ISBN != nil && strings.TrimSpace(*req.ISBN) == "" {
		return errors.New("isbn cannot be empty")
	}

	if req.Rating != nil {
		if *req.Rating < 1 || *req.Rating > 10 {
			return errors.New("rating must be between 1 and 10")
		}
	}

	return nil
}
