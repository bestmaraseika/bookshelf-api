package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	// Создаём HTTP-запрос
	req := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)

	recorder := httptest.NewRecorder()

	Health(recorder, req)

	resp := recorder.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			resp.StatusCode,
		)
	}

	contentType := resp.Header.Get("Content-Type")

	if contentType != "application/json" {
		t.Errorf(
			"expected application/json, got %s",
			contentType,
		)
	}

	expectedBody := "{\"status\":\"ok\"}\n"
	actualBody := recorder.Body.String()

	if actualBody != expectedBody {
		t.Errorf(
			"expected body %q, got %q",
			expectedBody,
			actualBody,
		)
	}
}
