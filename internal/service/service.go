package service

import (
	"github.com/bookshelf/monolith/internal/repository"
)

type Service struct {
	UserService *UserService
	BookService *BookService
}

func New(repos *repository.Repository, jwtSecret string) *Service {
	return &Service{
		UserService: NewUserService(repos.UserRepository, jwtSecret),
		BookService: NewBookService(repos.BookRepository, repos.UserRepository),
	}
}
