package service

import (
	"bookstore_api/internal/core/domain/books"
	"bookstore_api/internal/port"
	"context"
	"encoding/base64"
	"errors"
	"os"
)

var (
	KeyError      = errors.New("aes key is not set")
	DecodingError = errors.New("failed to decode")
)

type BookService struct {
	aesKey   []byte
	bookRepo port.BookRepository
}

func NewBookService(bookRepo port.BookRepository) (*BookService, error) {
	encodedKey := os.Getenv("AES_KEY")
	if encodedKey == "" {
		return nil, KeyError
	}

	key, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, DecodingError
	}

	return &BookService{
		aesKey:   key,
		bookRepo: bookRepo,
	}, nil
}

func (s *BookService) CreateBook(ctx context.Context, book *books.Book) (*books.Book, error) {
	err := book.EncryptCover(s.aesKey)
	if err != nil {
		return nil, err
	}

	return s.bookRepo.Create(ctx, book)
}

func (s *BookService) GetBookById(ctx context.Context, id int64) (*books.Book, error) {
	book, err := s.bookRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	err = book.DecryptCover(s.aesKey)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (s *BookService) GetAllBooks(ctx context.Context, page int64) ([]*books.Book, error) {
	allBook, err := s.bookRepo.GetAll(ctx, page)
	if err != nil {
		return nil, err
	}

	for _, book := range allBook {
		err = book.DecryptCover(s.aesKey)
		if err != nil {
			return nil, err
		}
	}

	return allBook, nil
}

func (s *BookService) UpdateBook(ctx context.Context, id int64, book *books.Book) (*books.Book, error) {
	// Fetch the existing book by its ID
	existingBook, err := s.bookRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	// Decrypt the cover image
	err = existingBook.DecryptCover(s.aesKey)
	if err != nil {
		return nil, err
	}
	{
		// Update fields if they have been modified
		if book.Title.Get() != "" {
			existingBook.UpdateTitle(book.Title.Get())
		}
		if book.CoverImage.Get() != "" {
			existingBook.CoverImage = book.CoverImage
		}
		if book.Synopsis.Get() != "" {
			existingBook.Synopsis = book.Synopsis
		}
		price := book.Price.Get()
		if price != 0 {
			err = existingBook.UpdatePrice(price)
			if err != nil {
				return nil, err
			}
		}
		stock := book.Stock.Get()
		if stock != 0 {
			err = existingBook.UpdateStock(stock)
			if err != nil {
				return nil, err
			}
		}
	}
	// Apply the update to the book
	existingBook.MarkUpdated()

	// Encrypt the updated cover image
	err = existingBook.EncryptCover(s.aesKey)
	if err != nil {
		return nil, err
	}

	// Save the updated book to the repository
	updatedBook, err := s.bookRepo.Update(ctx, existingBook)
	if err != nil {
		return nil, err
	}

	// Decrypt the cover image of the updated book before returning
	err = updatedBook.DecryptCover(s.aesKey)
	if err != nil {
		return nil, err
	}

	return updatedBook, nil
}

func (s *BookService) DeleteBook(ctx context.Context, id int64) error {
	return s.bookRepo.Delete(ctx, id)
}
