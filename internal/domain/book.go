package domain

import (
	"database/sql"
	"time"
)

type Book struct {
	ID            string          `json:"id" db:"id"`
	Title         string          `json:"title" db:"title"`
	Author        string          `json:"author" db:"author"`
	CreatedBy     string          `json:"created_by" db:"created_by"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at" db:"updated_at"`
	Description   sql.NullString  `json:"description" db:"description"`
	ISBN          sql.NullString  `json:"isbn" db:"isbn"`
	PublishedYear sql.NullInt32   `json:"published_year" db:"published_year"`
	AverageRating sql.NullFloat64 `json:"-" db:"average_rating"`
	ReviewsCount  int             `json:"reviews_count" db:"reviews_count"`
}

func (b *Book) ToResponse() BookResponse {
	book := BookResponse{
		ID:        b.ID,
		Title:     b.Title,
		Author:    b.Author,
		CreatedBy: b.CreatedBy,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}

	_ = b.Description.Scan(&book.Description)
	_ = b.ISBN.Scan(&book.ISBN)
	_ = b.PublishedYear.Scan(&book.PublishedYear)

	return book
}

type BookResponse struct {
	ID            string       `json:"id"`
	Title         string       `json:"title"`
	Author        string       `json:"author"`
	CreatedBy     string       `json:"created_by"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	Description   *string      `json:"description"`
	ISBN          *string      `json:"isbn"`
	PublishedYear *int         `json:"published_year"`
	AverageRating *float64     `json:"-"`
	ReviewsCount  int          `json:"reviews_count"`
	Creator       *UserSummary `json:"creator,omitempty"`
}

type CreateBookRequest struct {
	Title         string  `json:"title"`
	Author        string  `json:"author"`
	Description   *string `json:"description,omitempty"`
	ISBN          *int    `json:"isbn,omitempty"`
	PublishedYear *int    `json:"published_year,omitempty"`
}

type UpdateBookRequest struct {
	Title         *string `json:"title,omitempty"`
	Author        *string `json:"author,omitempty"`
	Description   *string `json:"description,omitempty"`
	ISBN          *int    `json:"isbn,omitempty"`
	PublishedYear *int    `json:"published_year,omitempty"`
}

type BookFilter struct {
	Search string `json:"search"`
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

type BookListResponse struct {
	Data       []BookResponse `json:"data"`
	Pagination Pagination     `json:"pagination"`
}
