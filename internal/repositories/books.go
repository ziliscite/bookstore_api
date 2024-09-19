package repositories

import (
	"bookstore_api/models"
	"context"
	"fmt"
)

type BookRepository interface {
	Create(ctx context.Context, book *models.Book) (*models.Book, error)
	GetById(ctx context.Context, id int) (*models.Book, error)
	GetAll(ctx context.Context) ([]*models.Book, error)
	Update(ctx context.Context, book *models.Book) (*models.Book, error)
	Delete(ctx context.Context, id int) error
}

func (r *Repository) Create(ctx context.Context, book *models.Book) (*models.Book, error) {
	// Execute query
	res, err := r.Db.NamedExecContext(ctx, "INSERT INTO books (title, slug, cover_image, synopsis, price, stock) VALUES (:title, :slug, :cover_image, :synopsis, :price, :stock)", book)
	if err != nil {
		return nil, fmt.Errorf("error while inserting book: %v", err)
	}

	// Getting the last inserted id (we're not manually inserting it up there)
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last inserted id")
	}

	// Adding the id
	book.ID = int(id)

	// Returning the book (reference) with id
	return book, nil
}

func (r *Repository) GetById(ctx context.Context, id int) (*models.Book, error) {
	var book models.Book
	err := r.Db.GetContext(ctx, &book, "SELECT * FROM books WHERE id=?", id)
	if err != nil {
		return nil, fmt.Errorf("error getting book: %v", err)
	}

	return &book, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]*models.Book, error) {
	var books []*models.Book
	err := r.Db.SelectContext(ctx, &books, "SELECT * FROM books")
	if err != nil {
		return nil, fmt.Errorf("error getting books: %v", err)
	}

	return books, nil
}

func (r *Repository) Update(ctx context.Context, book *models.Book) (*models.Book, error) {
	_, err := r.Db.NamedExecContext(ctx, "UPDATE books SET title=:title, slug=:slug, cover_image=:cover_image, synopsis=:synopsis, price=:price, stock=:stock WHERE id=:id", book)
	if err != nil {
		return nil, fmt.Errorf("error while updating book: %v", err)
	}

	return book, err
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	_, err := r.Db.ExecContext(ctx, "DELETE FROM books WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("error deleting book: %v", err)
	}

	return nil
}
