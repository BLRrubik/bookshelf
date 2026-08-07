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
	createUserQuery = `
INSERT INTO users (id, username, email, password_hash)
VALUES ($1, $2, $3, $4)
ON CONFLICT DO NOTHING;
`
	getUserByIDQuery = `
SELECT id, username, email, password_hash
FROM users
WHERE id = $1;
`
	getUserByUsernameQuery = `
SELECT id, username, email, password_hash
FROM users
WHERE username = $1;
`
	getUserByEmailQuery = `
SELECT id, username, email, password_hash
FROM users
WHERE email = $1;
`
	updateUserQuery = `
UPDATE users
SET username = $1, email = $2, password_hash = $3
WHERE id = $4;
`

	existsUserByUsernameQuery = `
SELECT EXISTS(SELECT id FROM users WHERE username = $1);
`

	existsUserByEmailQuery = `
SELECT EXISTS(SELECT id FROM users WHERE email = $1);
`
	getUsersByIDsQuery = `
SELECT id, username, email, password_hash
FROM users
WHERE id IN ($1);
`
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) Create(ctx context.Context, user *domain.User) error {
	user.ID = uuid.NewString()

	_, err := ur.db.ExecContext(ctx, createUserQuery, user.ID, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := ur.db.SelectContext(ctx, &user, getUserByIDQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := ur.db.SelectContext(ctx, &user, getUserByUsernameQuery, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := ur.db.SelectContext(ctx, &user, getUserByEmailQuery, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) Update(ctx context.Context, user *domain.User) error {
	_, err := ur.db.ExecContext(ctx, updateUserQuery, user.ID, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) UsernameExists(ctx context.Context, username string) bool {
	_, err := ur.db.ExecContext(ctx, existsUserByUsernameQuery, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false
		}
	}

	return true
}

func (ur *UserRepository) EmailExists(ctx context.Context, email string) bool {
	_, err := ur.db.ExecContext(ctx, existsUserByEmailQuery, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false
		}
	}

	return true
}

func (ur *UserRepository) GetByIDs(ctx context.Context, ids []string) (map[string]*domain.User, error) {
	result := make(map[string]*domain.User)

	return result, nil
}
