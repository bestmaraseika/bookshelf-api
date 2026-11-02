package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const maxBodySize = 1 << 20

func decodeBookRequest(
	w http.ResponseWriter,
	r *http.Request,
) (createBookRequest, error) {
	var req createBookRequest

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return req, fmt.Errorf("read request body: %w", err)
	}

	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return req, errors.New("request body is empty")
	}

	if trimmed[0] != '{' {
		return req, errors.New("expected a JSON object")
	}

	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		return req, fmt.Errorf("decode JSON: %w", err)
	}

	var extra any

	if err := decoder.Decode(&extra); err != io.EOF {
		return req, errors.New("request body must contain only one JSON object")
	}

	return req, nil
}
