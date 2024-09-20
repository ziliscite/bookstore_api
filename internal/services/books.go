package services

import (
	"bookstore_api/internal/repositories"
	"bookstore_api/models"
	"bookstore_api/tools"
	"context"
	"time"
)

type BookService struct {
	bookRepo repositories.IBookRepository
}

func NewBookService(bookRepo repositories.IBookRepository) *BookService {
	return &BookService{
		bookRepo: bookRepo,
	}
}

func (s *BookService) CreateBook(ctx context.Context, book *models.Book) (*models.Book, error) {
	book.Slug = tools.Slugify(book.Title)

	return s.bookRepo.Create(ctx, book)
}

func (s *BookService) GetBookById(ctx context.Context, id int) (*models.Book, error) {
	return s.bookRepo.GetById(ctx, int64(id))
}

func (s *BookService) GetAllBooks(ctx context.Context) ([]*models.Book, error) {
	return s.bookRepo.GetAll(ctx)
}

func (s *BookService) UpdateBook(ctx context.Context, book *models.Book) (*models.Book, error) {
	book.Slug = tools.Slugify(book.Title)

	t := time.Now().UTC()
	book.UpdatedAt = &t

	return s.bookRepo.Update(ctx, book)
}

func (s *BookService) DeleteBook(ctx context.Context, id int) error {
	return s.bookRepo.Delete(ctx, int64(id))
}
