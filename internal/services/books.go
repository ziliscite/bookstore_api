package services

import (
	"bookstore_api/internal/repositories"
	"bookstore_api/models"
	"context"
)

type BookService struct {
	bookRepo repositories.BookRepository
}

func NewBookService(bookRepo repositories.BookRepository) *BookService {
	return &BookService{bookRepo: bookRepo}
}

func (s *BookService) CreateBook(ctx context.Context, book *models.Book) (*models.Book, error) {
	return s.bookRepo.Create(ctx, book)
}

func (s *BookService) GetBookById(ctx context.Context, id int) (*models.Book, error) {
	return s.bookRepo.GetById(ctx, id)
}

func (s *BookService) GetAllBooks(ctx context.Context) ([]*models.Book, error) {
	return s.bookRepo.GetAll(ctx)
}

func (s *BookService) UpdateBook(ctx context.Context, book *models.Book) (*models.Book, error) {
	return s.bookRepo.Update(ctx, book)
}

func (s *BookService) DeleteBook(ctx context.Context, id int) error {
	return s.bookRepo.Delete(ctx, id)
}
