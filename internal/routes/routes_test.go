package routes

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestBookYearOutOfRange(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut} {
		for _, year := range []string{"2147483648", "-2147483649"} {
			t.Run(method+"/"+year, func(t *testing.T) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				defer db.Close()
				path := "/books"
				if method == http.MethodPut {
					path += "/1"
				}
				body := `{"title":"Dune","author":"Frank Herbert","year":` + year + `}`
				response := httptest.NewRecorder()
				NewRouter(db).ServeHTTP(response, httptest.NewRequest(method, path, strings.NewReader(body)))
				assertJSONError(t, response, http.StatusBadRequest)
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}

func TestBookYearMaxInteger(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut} {
		t.Run(method, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			now := time.Now()
			path, status := "/books", http.StatusCreated
			if method == http.MethodPost {
				mock.ExpectQuery("INSERT INTO books").WithArgs("Dune", "Frank Herbert", int64(math.MaxInt32), nil, false, nil).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(1, now, now))
			} else {
				path, status = "/books/1", http.StatusOK
				mock.ExpectQuery("UPDATE books").WithArgs("Dune", "Frank Herbert", int64(math.MaxInt32), nil, false, nil, int64(1)).
					WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).AddRow(now, now))
			}
			response := httptest.NewRecorder()
			NewRouter(db).ServeHTTP(response, httptest.NewRequest(method, path,
				strings.NewReader(`{"title":"Dune","author":"Frank Herbert","year":2147483647}`)))
			if response.Code != status {
				t.Fatalf("status %d, want %d: %s", response.Code, status, response.Body.String())
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestPaginationOffsetOverflow(t *testing.T) {
	if strconv.IntSize < 64 {
		t.Skip("offset overflow requires 64-bit int")
	}
	for _, query := range []string{
		"page=9223372036854775807&limit=100",
		"page=92233720368547759&limit=100",
	} {
		t.Run(query, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			response := httptest.NewRecorder()
			NewRouter(db).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/books?"+query, nil))
			assertJSONError(t, response, http.StatusBadRequest)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestPaginationLargestAcceptedOffset(t *testing.T) {
	if strconv.IntSize < 64 {
		t.Skip("boundary requires 64-bit int")
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	// Largest page accepted by the existing repository guard for limit=100.
	const page = (math.MaxInt64-100)/100 + 1
	const offset = (page - 1) * 100
	mock.ExpectQuery("SELECT").WithArgs(100, int64(offset)).WillReturnRows(sqlmock.NewRows(
		[]string{"id", "title", "author", "year", "isbn", "out_of_stock", "rating", "created_at", "updated_at"}))
	response := httptest.NewRecorder()
	NewRouter(db).ServeHTTP(response, httptest.NewRequest(http.MethodGet,
		"/books?page="+strconv.FormatInt(page, 10)+"&limit=100", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("status %d: %s", response.Code, response.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestRoutingJSONErrors(t *testing.T) {
	tests := []struct {
		name, method, path string
		status             int
		allow              []string
	}{
		{"unknown path", http.MethodGet, "/missing", http.StatusNotFound, nil},
		{"unknown path with unsupported method", http.MethodPatch, "/missing", http.StatusNotFound, nil},
		{"collection", http.MethodPatch, "/books", http.StatusMethodNotAllowed, []string{"GET", "POST"}},
		{"book", http.MethodPost, "/books/1", http.StatusMethodNotAllowed, []string{"DELETE", "GET", "PUT"}},
		{"recommend", http.MethodPatch, "/books/recommend", http.StatusMethodNotAllowed, []string{"GET", "PUT", "DELETE"}},
		{"stock", http.MethodGet, "/books/1/mark-out-of-stock", http.StatusMethodNotAllowed, []string{"POST"}},
		{"health", http.MethodPost, "/health", http.StatusMethodNotAllowed, []string{"GET"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			NewRouter(nil).ServeHTTP(response, httptest.NewRequest(tt.method, tt.path, nil))
			assertJSONError(t, response, tt.status)
			var got []string
			for _, value := range response.Header().Values("Allow") {
				for _, method := range strings.Split(value, ",") {
					if method = strings.TrimSpace(method); method != "" {
						got = append(got, method)
					}
				}
			}
			sort.Strings(got)
			sort.Strings(tt.allow)
			if strings.Join(got, ",") != strings.Join(tt.allow, ",") {
				t.Fatalf("Allow %v, want %v", got, tt.allow)
			}
		})
	}
}

func assertJSONError(t *testing.T, response *httptest.ResponseRecorder, status int) {
	t.Helper()
	if response.Code != status {
		t.Fatalf("status %d, want %d: %s", response.Code, status, response.Body.String())
	}
	if got := response.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type %q", got)
	}
	var body map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if body["error"] == "" {
		t.Errorf("missing error: %s", response.Body.String())
	}
}
