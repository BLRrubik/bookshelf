package service

import (
	"context"
	"errors"

	"github.com/bookshelf/monolith/internal/domain"
	"github.com/bookshelf/monolith/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrUsernameExists     = errors.New("username already exists")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidUsername    = errors.New("invalid username")
	ErrInvalidEmail       = errors.New("invalid email")
)

type UserService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewUserService(userRepo *repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *UserService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error) {
	if err := s.validateRegister(ctx, req); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
	}

	if err = s.userRepo.Create(ctx, &user); err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		User:        user.ToPublic(),
		AccessToken: "",
		TokenType:   "Bearer",
		ExpiresIn:   0,
	}, nil
}

func (s *UserService) validateRegister(ctx context.Context, req domain.RegisterRequest) error {
	if len(req.Email) == 0 {
		return ErrInvalidEmail
	}

	if len(req.Username) < 3 {
		return ErrInvalidUsername
	}

	if len(req.Password) < 8 {
		return ErrInvalidPassword
	}

	if s.userRepo.EmailExists(ctx, req.Email) {
		return ErrUserExists
	}

	if s.userRepo.UsernameExists(ctx, req.Username) {
		return ErrUsernameExists
	}

	return nil
}
