package services

import (
	"bookstore_api/internal/repositories"
	"bookstore_api/models"
	"bookstore_api/tools"
	"context"
	"errors"
	"regexp"
	"strconv"
	"time"
)

type BookService struct {
	*Service
	bookRepo repositories.IBookRepository
}

func NewBookService(service *Service, bookRepo repositories.IBookRepository) *BookService {
	return &BookService{
		Service:  service,
		bookRepo: bookRepo,
	}
}

func (s *BookService) CreateBook(ctx context.Context, book *models.Book) (*models.Book, error) {
	book.Slug = tools.Slugify(book.Title)

	err := s.encryptBook(book)
	if err != nil {
		return nil, err
	}

	return s.bookRepo.Create(ctx, book)
}

func (s *BookService) GetBookById(ctx context.Context, idStr string) (*models.Book, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, errors.New("invalid Id")
	}

	book, err := s.bookRepo.GetById(ctx, int64(id))
	if err != nil {
		return nil, err
	}

	err = s.decryptBook(book)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (s *BookService) GetAllBooks(ctx context.Context) ([]*models.Book, error) {
	page, exist := ctx.Value("page").(string)
	if !exist || page == "" {
		page = "1"
	}

	// Define a regex pattern to allow only positive integers
	re := regexp.MustCompile(`^[1-9]\d*$`)
	if !re.MatchString(page) {
		return nil, errors.New("page not valid")
	}

	pageNum, err := strconv.Atoi(page)
	if err != nil {
		return nil, errors.New("error converting string to integer")
	}

	books, err := s.bookRepo.GetAll(ctx, pageNum)
	if err != nil {
		return nil, err
	}

	for _, book := range books {
		err = s.decryptBook(book)
		if err != nil {
			return nil, err
		}
	}

	return books, nil
}

func (s *BookService) UpdateBook(ctx context.Context, book *models.Book) (*models.Book, error) {
	book.Slug = tools.Slugify(book.Title)

	t := time.Now()
	book.UpdatedAt = &t

	err := s.encryptBook(book)
	if err != nil {
		return nil, err
	}

	// Updates book
	updatedBook, err := s.bookRepo.Update(ctx, book)
	if err != nil {
		return nil, err
	}

	// Decrypt the URL for return value
	err = s.decryptBook(updatedBook)
	if err != nil {
		return nil, err
	}

	return updatedBook, nil
}

func (s *BookService) DeleteBook(ctx context.Context, idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return errors.New("invalid Id")
	}

	return s.bookRepo.Delete(ctx, int64(id))
}

// A generic name, since we'll encrypt/decrypt whatever might be necessary in the future, not just the image's URL
// Since the book might later be JOIN'ed with other data?
func (s *BookService) encryptBook(book *models.Book) error {
	encryptedURL, err := tools.Encrypt([]byte(book.CoverImage), s.AESKey)
	if err != nil {
		return err
	}

	book.CoverImage = encryptedURL
	return nil
}

func (s *BookService) decryptBook(book *models.Book) error {
	decryptedURL, err := tools.Decrypt(book.CoverImage, s.AESKey)
	if err != nil {
		return err
	}

	book.CoverImage = string(decryptedURL)
	return nil
}
