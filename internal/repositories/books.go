package repositories

import (
	"bookstore_api/models"
	"context"
	"errors"
	"fmt"
)

type BookRepository struct {
	*Repository
}

func NewBookRepository(repository *Repository) *BookRepository {
	return &BookRepository{
		repository,
	}
}

type IBookRepository interface {
	Create(ctx context.Context, book *models.Book) (*models.Book, error)
	GetById(ctx context.Context, id int64) (*models.Book, error)
	GetAll(ctx context.Context, page int) ([]*models.Book, error)
	Update(ctx context.Context, book *models.Book) (*models.Book, error)
	Delete(ctx context.Context, id int64) error
}

func (repo *BookRepository) Create(ctx context.Context, book *models.Book) (*models.Book, error) {
	query := `
		INSERT INTO books (title, slug, cover_image, synopsis, price, stock) 
		VALUES (:title, :slug, :cover_image, :synopsis, :price, :stock) 
		ON CONFLICT (title)
		DO NOTHING
		RETURNING *
    `

	// Prepare the named query
	stmt, err := repo.Db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparing query: %v", err)
	}

	// Execute the query and scan the result into `book` directly
	err = stmt.GetContext(ctx, book, book)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("error while inserting book: book title already exists")
		}
		return nil, fmt.Errorf("error while inserting book: %v", err)
	}

	return book, nil
}

func (repo *BookRepository) GetById(ctx context.Context, id int64) (*models.Book, error) {
	book := &models.Book{}
	err := repo.Db.GetContext(ctx, book, "SELECT * FROM books WHERE id=$1", id)
	if err != nil {
		return nil, fmt.Errorf("error getting book: %v", err)
	}

	return book, nil
}

func (repo *BookRepository) GetAll(ctx context.Context, page int) ([]*models.Book, error) {
	limit := 20
	offset := limit * (page - 1)

	query := `
		SELECT * 
		FROM books 
		LIMIT $1 
		OFFSET $2
	`

	var books []*models.Book
	err := repo.Db.SelectContext(ctx, &books, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting books: %v", err)
	}

	return books, nil
}

func (repo *BookRepository) Update(ctx context.Context, book *models.Book) (*models.Book, error) {
	query := `
		WITH title_conflict AS (
			SELECT id FROM books WHERE title = :title AND NOT id = :id
		)
		UPDATE books
		SET title=:title, slug=:slug, cover_image=:cover_image, synopsis=:synopsis, price=:price, stock=:stock, updated_at=:updated_at 
		WHERE id=:id AND NOT EXISTS (SELECT 1 FROM title_conflict)
		RETURNING *
	`

	// Use NamedQueryRowContext for queries that return a single row.
	stmt, err := repo.Db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparing query: %v", err)
	}

	// Execute the query and map the result to updatedBook.
	updatedBook := &models.Book{}
	err = stmt.GetContext(ctx, updatedBook, book)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("error while updating book: book title already exists")
		}
		return nil, fmt.Errorf("error while updating book: %v", err)
	}

	return updatedBook, nil
}

func (repo *BookRepository) Delete(ctx context.Context, id int64) error {
	_, err := repo.Db.ExecContext(ctx, "DELETE FROM books WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("error deleting book: %v", err)
	}

	return nil
}
