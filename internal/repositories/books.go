package repositories

import (
	"bookstore_api/models"
	"context"
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
	GetAll(ctx context.Context) ([]*models.Book, error)
	Update(ctx context.Context, book *models.Book) (*models.Book, error)
	Delete(ctx context.Context, id int64) error
}

func (r *BookRepository) Create(ctx context.Context, book *models.Book) (*models.Book, error) {
	query := `
        INSERT INTO books (title, slug, cover_image, synopsis, price, stock) 
        VALUES (:title, :slug, :cover_image, :synopsis, :price, :stock) 
        RETURNING *
    `

	// Prepare the named query
	stmt, err := r.Db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparing query: %v", err)
	}

	// Execute the query and scan the result into `book` directly
	err = stmt.GetContext(ctx, book, book)
	if err != nil {
		return nil, fmt.Errorf("error while inserting book: %v", err)
	}

	return book, nil
}
func (r *BookRepository) GetById(ctx context.Context, id int64) (*models.Book, error) {
	var book models.Book
	err := r.Db.GetContext(ctx, &book, "SELECT * FROM books WHERE id=$1", id)
	if err != nil {
		return nil, fmt.Errorf("error getting book: %v", err)
	}

	return &book, nil
}

func (r *BookRepository) GetAll(ctx context.Context) ([]*models.Book, error) {
	var books []*models.Book
	err := r.Db.SelectContext(ctx, &books, "SELECT * FROM books")
	if err != nil {
		return nil, fmt.Errorf("error getting books: %v", err)
	}

	return books, nil
}

func (r *BookRepository) Update(ctx context.Context, book *models.Book) (*models.Book, error) {
	query := `
		UPDATE books 
		SET title=:title, slug=:slug, cover_image=:cover_image, synopsis=:synopsis, price=:price, stock=:stock 
		WHERE id=:id
		RETURNING *
	`

	// Use NamedQueryRowContext for queries that return a single row.
	stmt, err := r.Db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparing query: %v", err)
	}

	// Execute the query and map the result to updatedBook.
	updatedBook := &models.Book{}
	err = stmt.GetContext(ctx, updatedBook, book)
	if err != nil {
		return nil, fmt.Errorf("error while updating book: %v", err)
	}

	return updatedBook, nil
}

func (r *BookRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.Db.ExecContext(ctx, "DELETE FROM books WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("error deleting book: %v", err)
	}

	return nil
}
