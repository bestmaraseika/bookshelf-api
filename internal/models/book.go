package models

import "time"

type Book struct {
	ID         uint64    `json:"id"`
	Title      string    `json:"title"`
	Author     string    `json:"author"`
	Year       int       `json:"year"`
	ISBN       *string   `json:"isbn"`
	OutOfStock bool      `json:"out_of_stock"`
	Rating     *int      `json:"rating"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
