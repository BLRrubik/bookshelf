package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bookshelf/monolith/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	createBookQuery = `
INSERT INTO books 
(id, title, author, description, isbn, published_year, created_by)
VALUES($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT DO NOTHING;
`
	getBookByIDQuery = `
SELECT 
	id, title, author, 
	description, isbn, published_year,
	created_by, created_at, updated_at 
FROM books
WHERE id = $1
`
	listBooksQuery = `
SELECT 
	id, title, author, 
	description, isbn, published_year,
	created_by, created_at, updated_at
FROM books
WHERE title ILIKE $1
ORDER BY $2 $3
LIMIT $4 OFFSET $5
`
	updateBookQuery = `
UPDATE books
SET title=$1, author=$2, description=$3, isbn=$4, published_year=$5
WHERE id = $6
`
	deleteBookQuery = `
DELETE FROM books
WHERE id = $1
`
)

type BookRepository struct {
	db *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) *BookRepository {
	return &BookRepository{
		db: db,
	}
}

func (br *BookRepository) Create(ctx context.Context, book *domain.Book) error {
	book.ID = uuid.NewString()

	_, err := br.db.ExecContext(
		ctx,
		createBookQuery,
		book.ID,
		book.Title,
		book.Author,
		book.Description,
		book.ISBN,
		book.PublishedYear,
		book.CreatedBy,
	)
	if err != nil {
		return err
	}

	return nil
}

func (br *BookRepository) GetByID(ctx context.Context, id string) (*domain.Book, error) {
	var book domain.Book
	err := br.db.SelectContext(ctx, &book, getBookByIDQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBookNotFound
		}

		return nil, err
	}

	return &book, nil
}

func (br *BookRepository) List(ctx context.Context, filter domain.BookFilter) ([]domain.Book, int, error) {
	var books []domain.Book
	var count int

	offset := (filter.Page - 1) * filter.Limit

	err := br.db.SelectContext(
		ctx,
		&books,
		listBooksQuery,
		filter.Search,
		filter.Order,
		filter.Sort,
		filter.Limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}

	err = br.db.SelectContext(ctx, &count, "SELECT COUNT(*) FROM books")
	if err != nil {
		return nil, 0, err
	}

	return books, count, nil
}

func (br *BookRepository) Update(ctx context.Context, book *domain.Book) error {
	_, err := br.db.ExecContext(
		ctx,
		updateBookQuery,
		book.Title,
		book.Author,
		book.Description,
		book.ISBN,
		book.PublishedYear,
		book.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (br *BookRepository) Delete(ctx context.Context, id string) error {
	_, err := br.db.ExecContext(ctx, deleteBookQuery, id)
	if err != nil {
		return err
	}

	return nil
}
