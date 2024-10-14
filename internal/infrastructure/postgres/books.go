package postgres

import (
	"bookstore_api/internal/core/domain/books"
	"context"
	"fmt"
	"time"
)

type BookRepository struct {
	*Database
}

func NewBookRepository(db *Database) *BookRepository {
	return &BookRepository{
		Database: db,
	}
}

type dbBookDTO struct {
	ID    int64  `db:"id"`
	Title string `db:"title"`

	Slug       string `db:"slug"`
	CoverImage string `db:"cover_image"`
	Synopsis   string `db:"synopsis"`

	Price float64 `db:"price"`
	Stock int64   `db:"stock"`

	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func (d *dbBookDTO) newBook(id int64) (*books.Book, error) {
	book, err := books.NewBook(d.Title, d.CoverImage, d.Synopsis, d.Price, d.Stock)
	if err != nil {
		return nil, err
	}

	if d.CreatedAt != nil {
		book.CreatedAt = d.CreatedAt
	}

	err = book.Create(id)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func newBookDTO(book *books.Book) *dbBookDTO {
	return &dbBookDTO{
		ID:         book.ID.Get(),
		Title:      book.Title.Get(),
		Slug:       book.Slug.Get(),
		CoverImage: book.CoverImage.Get(),
		Synopsis:   book.Synopsis.Get(),
		Price:      book.Price.Get(),
		Stock:      book.Stock.Get(),
		CreatedAt:  book.CreatedAt,
		UpdatedAt:  book.UpdatedAt,
	}
}

func (r *BookRepository) Create(ctx context.Context, book *books.Book) (*books.Book, error) {
	bookDTO := newBookDTO(book)

	query := `
		INSERT INTO books (title, slug, cover_image, synopsis, price, stock) 
		VALUES (:title, :slug, :cover_image, :synopsis, :price, :stock) 
		ON CONFLICT (title)
		DO NOTHING
		RETURNING *
    `

	// Prepare the named query
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparing query: %v", err)
	}

	// Execute the query and scan the result into `bookDTO` directly
	err = stmt.GetContext(ctx, bookDTO, bookDTO)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("error while inserting book: book title already exists")
		}
		return nil, fmt.Errorf("error while inserting book: %v", err)
	}

	// This part's kinda sus
	newBook, err := bookDTO.newBook(bookDTO.ID)
	if err != nil {
		return nil, err
	}
	//

	return newBook, nil
}

func (r *BookRepository) GetById(ctx context.Context, id int64) (*books.Book, error) {
	book := &dbBookDTO{}
	err := r.db.GetContext(ctx, book, "SELECT * FROM books WHERE id=$1", id)
	if err != nil {
		return nil, fmt.Errorf("error getting book: %v", err)
	}

	newBook, err := book.newBook(id)
	if err != nil {
		return nil, err
	}

	return newBook, nil
}

func (r *BookRepository) GetAll(ctx context.Context, page int64) ([]*books.Book, error) {
	limit := 20
	offset := limit * (int(page) - 1)

	query := `
		SELECT * 
		FROM books 
		LIMIT $1 
		OFFSET $2
	`

	var booksDTO []*dbBookDTO
	err := r.db.SelectContext(ctx, &booksDTO, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting books: %v", err)
	}

	var allBooks []*books.Book
	for _, book := range booksDTO {
		newBook, err := book.newBook(book.ID)
		if err != nil {
			return nil, err
		}
		allBooks = append(allBooks, newBook)
	}

	return allBooks, nil
}

func (r *BookRepository) Update(ctx context.Context, book *books.Book) (*books.Book, error) {
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
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparing query: %v", err)
	}

	// Execute the query and map the result to updatedBook.
	bookDTO := newBookDTO(book)
	err = stmt.GetContext(ctx, bookDTO, bookDTO)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("error while updating book: book title already exists")
		}
		return nil, fmt.Errorf("error while updating book: %v", err)
	}

	newBook, err := bookDTO.newBook(book.ID.Get())
	if err != nil {
		return nil, err
	}

	return newBook, nil
}

func (r *BookRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM books WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("error deleting book: %v", err)
	}

	return nil
}
