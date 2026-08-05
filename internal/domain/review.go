package domain

import (
	"database/sql"
	"time"
)

type Review struct {
	ID        string         `json:"id" db:"id"`
	BookID    string         `json:"book_id" db:"book_id"`
	UserID    string         `json:"user_id" db:"user_id"`
	Rating    int            `json:"rating" db:"rating"`
	Title     sql.NullString `json:"title" db:"title"`
	Content   string         `json:"content" db:"content"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

func (r *Review) ToResponse(user *User) ReviewResponse {
	review := ReviewResponse{
		ID:        r.ID,
		BookID:    r.BookID,
		UserID:    r.UserID,
		Rating:    r.Rating,
		Content:   r.Content,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}

	if user != nil {
		review.User = user.ToSummary()
	}

	_ = r.Title.Scan(&review.Title)

	return review
}

type ReviewResponse struct {
	ID        string      `json:"id"`
	BookID    string      `json:"book_id"`
	UserID    string      `json:"user_id"`
	Rating    int         `json:"rating"`
	Title     *string     `json:"title"`
	Content   string      `json:"content"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	User      UserSummary `json:"user"`
}

type CreateReviewRequest struct {
	Rating  int     `json:"rating"`
	Content string  `json:"content"`
	Title   *string `json:"title,omitempty"`
}

type UpdateReviewRequest struct {
	Rating  *int    `json:"rating,omitempty"`
	Content *string `json:"content,omitempty"`
	Title   *string `json:"title,omitempty"`
}

type ReviewListResponse struct {
	Data       []ReviewResponse `json:"data"`
	Pagination Pagination       `json:"pagination"`
}
