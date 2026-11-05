package handlers

import "testing"

func TestValidateBookRequest_Valid(t *testing.T) {
	rating := 9

	req := createBookRequest{
		Title:  "Dune",
		Author: "Frank Herbert",
		Year:   1965,
		Rating: &rating,
	}

	err := validateBookRequest(req)

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateBookRequest_Table(t *testing.T) {
	validRating := 9
	invalidRating := 15
	emptyISBN := ""

	tests := []struct {
		name    string
		request createBookRequest
		wantErr bool
	}{
		{
			name: "valid book",
			request: createBookRequest{
				Title:  "Dune",
				Author: "Frank Herbert",
				Year:   1965,
				Rating: &validRating,
			},
			wantErr: false,
		},
		{
			name: "empty title",
			request: createBookRequest{
				Title:  "",
				Author: "Frank Herbert",
				Year:   1965,
			},
			wantErr: true,
		},
		{
			name: "title with spaces",
			request: createBookRequest{
				Title:  "   ",
				Author: "Frank Herbert",
				Year:   1965,
			},
			wantErr: true,
		},
		{
			name: "empty author",
			request: createBookRequest{
				Title:  "Dune",
				Author: "",
				Year:   1965,
			},
			wantErr: true,
		},
		{
			name: "negative year",
			request: createBookRequest{
				Title:  "Dune",
				Author: "Frank Herbert",
				Year:   -10,
			},
			wantErr: true,
		},
		{
			name: "invalid rating",
			request: createBookRequest{
				Title:  "Dune",
				Author: "Frank Herbert",
				Year:   1965,
				Rating: &invalidRating,
			},
			wantErr: true,
		},
		{
			name: "missing rating",
			request: createBookRequest{
				Title:  "Dune",
				Author: "Frank Herbert",
				Year:   1965,
				Rating: nil,
			},
			wantErr: false,
		},
		{
			name: "empty isbn",
			request: createBookRequest{
				Title:  "Dune",
				Author: "Frank Herbert",
				Year:   1965,
				ISBN:   &emptyISBN,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBookRequest(tt.request)

			gotErr := err != nil

			if gotErr != tt.wantErr {
				t.Errorf(
					"validateBookRequest() error = %v, wantErr = %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}
