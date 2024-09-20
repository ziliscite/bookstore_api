package tests

import (
	"bookstore_api/internal/repositories"
	"bookstore_api/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"testing"
)

// TODO: Test the services & tools

func NewDatabaseMock() (*sql.DB, sqlmock.Sqlmock, error) {
	// Mocking database
	conn, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual)) // A more sensitive query
	if err != nil {
		return nil, nil, err
	}

	return conn, mock, nil
}

func withDatabaseMock(t *testing.T, fn func(*sqlx.DB, sqlmock.Sqlmock)) {
	conn, mock, err := NewDatabaseMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer func(conn *sql.DB) {
		err = conn.Close()
		if err != nil {

		}
	}(conn)

	db := sqlx.NewDb(conn, "sqlmock")
	fn(db, mock)
}

func TestCreateBook(t *testing.T) {
	book := &models.Book{
		Title:      "Solo Leveling",
		Slug:       "solo-leveling",
		CoverImage: "https://asura.nacmcdn.com/wp-content/uploads/2022/03/solo-leveling.jpg",
		Synopsis:   "In a world where hunters, humans with magical abilities...",
		Price:      12.99,
		Stock:      517,
	}

	cases := []struct {
		name string
		test func(*testing.T, *repositories.BookRepository, sqlmock.Sqlmock)
	}{
		{
			name: "Success",
			test: func(t *testing.T, r *repositories.BookRepository, mock sqlmock.Sqlmock) {
				ctx := context.Background()
				bookId := 1

				// Change the query to match positional placeholders
				query := `
                    INSERT INTO books (title, slug, cover_image, synopsis, price, stock) 
                    VALUES (?, ?, ?, ?, ?, ?) 
                    RETURNING *
                `

				// Mocking the query with positional placeholders
				mock.ExpectPrepare(query)
				mock.ExpectQuery(query).
					WithArgs(book.Title, book.Slug, book.CoverImage, book.Synopsis, book.Price, book.Stock).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "slug", "cover_image", "synopsis", "price", "stock"}).
						AddRow(bookId, book.Title, book.Slug, book.CoverImage, book.Synopsis, book.Price, book.Stock))

				createdBook, err := r.Create(ctx, book)
				require.NoError(t, err)
				require.NotNil(t, createdBook)
				require.Equal(t, int64(bookId), createdBook.ID)

				// Ensure all expectations are met
				if err = mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			},
		},
		{
			name: "Exec Failure",
			test: func(t *testing.T, r *repositories.BookRepository, mock sqlmock.Sqlmock) {
				ctx := context.Background()

				query := `
                    INSERT INTO books (title, slug, cover_image, synopsis, price, stock) 
                    VALUES (?, ?, ?, ?, ?, ?) 
                    RETURNING *
                `

				// Expect failure during execution with positional placeholders
				mock.ExpectPrepare(query)
				mock.ExpectQuery(query).
					WithArgs(book.Title, book.Slug, book.CoverImage, book.Synopsis, book.Price, book.Stock).
					WillReturnError(fmt.Errorf("error while inserting book"))

				_, err := r.Create(ctx, book)
				require.Error(t, err)

				if err = mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			withDatabaseMock(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				mockRepo := repositories.NewBookRepository(repositories.NewRepository(db))
				c.test(t, mockRepo, mock)
			})
		})
	}
}

func TestGetById(t *testing.T) {
	book := &models.Book{
		ID:         1,
		Title:      "Solo Leveling",
		Slug:       "solo-leveling",
		CoverImage: "https://asura.nacmcdn.com/wp-content/uploads/2022/03/solo-leveling.jpg",
		Synopsis:   "In a world where hunters, humans with magical abilities...",
		Price:      12.99,
		Stock:      517,
	}

	cases := []struct {
		name string
		test func(*testing.T, *repositories.BookRepository, sqlmock.Sqlmock)
	}{
		{
			name: "Success",
			test: func(t *testing.T, r *repositories.BookRepository, mock sqlmock.Sqlmock) {
				ctx := context.Background()
				bookId := int64(1)

				// Mock the query for getting the book by ID
				query := `SELECT * FROM books WHERE id=$1`

				rows := sqlmock.NewRows([]string{"id", "title", "slug", "cover_image", "synopsis", "price", "stock"}).
					AddRow(book.ID, book.Title, book.Slug, book.CoverImage, book.Synopsis, book.Price, book.Stock)

				mock.ExpectQuery(query).WithArgs(bookId).WillReturnRows(rows)

				// Call the repository method
				retrievedBook, err := r.GetById(ctx, bookId)
				require.NoError(t, err)
				require.NotNil(t, retrievedBook)
				require.Equal(t, book.ID, retrievedBook.ID)
				require.Equal(t, book.Title, retrievedBook.Title)

				// Ensure all expectations are met
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			},
		},
		{
			name: "NotFound",
			test: func(t *testing.T, r *repositories.BookRepository, mock sqlmock.Sqlmock) {
				ctx := context.Background()
				bookId := int64(1)

				// Mock the query for getting the book by ID
				query := `SELECT * FROM books WHERE id=$1`

				mock.ExpectQuery(query).WithArgs(bookId).WillReturnError(sql.ErrNoRows)

				// Call the repository method
				retrievedBook, err := r.GetById(ctx, bookId)
				require.Error(t, err)
				require.Nil(t, retrievedBook)

				// Ensure all expectations are met
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			},
		},
		{
			name: "Query Error",
			test: func(t *testing.T, r *repositories.BookRepository, mock sqlmock.Sqlmock) {
				ctx := context.Background()
				bookId := int64(1)

				// Mock a query error
				query := `SELECT * FROM books WHERE id=$1`

				mock.ExpectQuery(query).WithArgs(bookId).WillReturnError(fmt.Errorf("query error"))

				// Call the repository method
				retrievedBook, err := r.GetById(ctx, bookId)
				require.Error(t, err)
				require.Nil(t, retrievedBook)

				// Ensure all expectations are met
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			withDatabaseMock(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				mockRepo := repositories.NewBookRepository(repositories.NewRepository(db))
				c.test(t, mockRepo, mock)
			})
		})
	}
}

func TestGetAll(t *testing.T) {
	books := []*models.Book{
		{
			ID:         1,
			Title:      "Omniscient Reader",
			Slug:       "omniscient-reader",
			CoverImage: "https://asura.nacmcdn.com/wp-content/uploads/2021/12/omniscient-reader.jpg",
			Synopsis:   "Dokja was an average office worker...",
			Price:      10.99,
			Stock:      432,
		},
		{
			ID:         2,
			Title:      "The Beginning After The End",
			Slug:       "the-beginning-after-the-end",
			CoverImage: "https://asura.nacmcdn.com/wp-content/uploads/2022/04/the-beginning-after-the-end.jpg",
			Synopsis:   "King Grey has unrivaled strength...",
			Price:      11.49,
			Stock:      289,
		},
	}

	cases := []struct {
		name string
		test func(*testing.T, *repositories.BookRepository, sqlmock.Sqlmock)
	}{
		{
			name: "Success",
			test: func(t *testing.T, r *repositories.BookRepository, mock sqlmock.Sqlmock) {
				ctx := context.Background()

				// Mock the query for getting all books
				query := `SELECT * FROM books`

				rows := sqlmock.NewRows([]string{"id", "title", "slug", "cover_image", "synopsis", "price", "stock"})

				for _, book := range books {
					rows.AddRow(book.ID, book.Title, book.Slug, book.CoverImage, book.Synopsis, book.Price, book.Stock)
				}

				mock.ExpectQuery(query).WillReturnRows(rows)

				// Call the repository method
				retrievedBooks, err := r.GetAll(ctx)
				require.NoError(t, err)
				require.NotNil(t, retrievedBooks)
				require.Len(t, retrievedBooks, 2)

				// Ensure all expectations are met
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			},
		},
		{
			name: "Query Error",
			test: func(t *testing.T, r *repositories.BookRepository, mock sqlmock.Sqlmock) {
				ctx := context.Background()

				// Mock a query error
				query := `SELECT * FROM books`

				mock.ExpectQuery(query).WillReturnError(fmt.Errorf("query error"))

				// Call the repository method
				retrievedBooks, err := r.GetAll(ctx)
				require.Error(t, err)
				require.Nil(t, retrievedBooks)

				// Ensure all expectations are met
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			withDatabaseMock(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				mockRepo := repositories.NewBookRepository(repositories.NewRepository(db))
				c.test(t, mockRepo, mock)
			})
		})
	}
}

func TestUpdateBook(t *testing.T) {
	book := &models.Book{
		Title:      "Solo Leveling",
		Slug:       "solo-leveling",
		CoverImage: "https://asura.nacmcdn.com/wp-content/uploads/2022/03/solo-leveling.jpg",
		Synopsis:   "In a world where hunters, humans with magical abilities...",
		Price:      12.99,
		Stock:      517,
		ID:         1, // Make sure ID is set here
	}

	postUpdateBook := &models.Book{
		Title:      "Solo Leveling",
		Slug:       "solo-leveling",
		CoverImage: "https://asura.nacmcdn.com/wp-content/uploads/2022/03/solo-leveling.jpg",
		Synopsis:   "In a world where hunters, humans with magical abilities...",
		Price:      15.99,
		Stock:      517,
		ID:         1, // Ensure the ID is set when updating the book
	}

	cases := []struct {
		name string
		test func(*testing.T, *repositories.BookRepository, sqlmock.Sqlmock)
	}{
		{
			name: "Success",
			test: func(t *testing.T, r *repositories.BookRepository, mock sqlmock.Sqlmock) {
				ctx := context.Background()
				bookId := 1

				// Mock the Insert Query for creating the book
				mock.ExpectPrepare("INSERT INTO books (title, slug, cover_image, synopsis, price, stock) VALUES (?, ?, ?, ?, ?, ?) RETURNING *")
				mock.ExpectQuery("INSERT INTO books (title, slug, cover_image, synopsis, price, stock) VALUES (?, ?, ?, ?, ?, ?) RETURNING *").
					WithArgs(book.Title, book.Slug, book.CoverImage, book.Synopsis, book.Price, book.Stock).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "slug", "cover_image", "synopsis", "price", "stock"}).
						AddRow(bookId, book.Title, book.Slug, book.CoverImage, book.Synopsis, book.Price, book.Stock))

				createdBook, err := r.Create(ctx, book)
				require.NoError(t, err)
				require.NotNil(t, createdBook)
				require.Equal(t, createdBook.ID, int64(bookId))
				require.Equal(t, createdBook, book)

				// Mock the Update Query
				mock.ExpectPrepare("UPDATE books SET title=?, slug=?, cover_image=?, synopsis=?, price=?, stock=? WHERE id=? RETURNING *")
				mock.ExpectQuery("UPDATE books SET title=?, slug=?, cover_image=?, synopsis=?, price=?, stock=? WHERE id=? RETURNING *").
					WithArgs(postUpdateBook.Title, postUpdateBook.Slug, postUpdateBook.CoverImage, postUpdateBook.Synopsis, postUpdateBook.Price, postUpdateBook.Stock, bookId).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "slug", "cover_image", "synopsis", "price", "stock"}).
						AddRow(bookId, postUpdateBook.Title, postUpdateBook.Slug, postUpdateBook.CoverImage, postUpdateBook.Synopsis, postUpdateBook.Price, postUpdateBook.Stock))

				updatedBook, err := r.Update(ctx, postUpdateBook)
				require.NoError(t, err)
				require.NotNil(t, updatedBook)
				require.Equal(t, updatedBook, postUpdateBook)

				// Check that all expectations were met
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			},
		},
		{
			name: "NamedExecContext Failure",
			test: func(t *testing.T, r *repositories.BookRepository, mock sqlmock.Sqlmock) {
				ctx := context.Background()

				// Mock the Update Query to fail
				mock.ExpectPrepare("UPDATE books SET title=?, slug=?, cover_image=?, synopsis=?, price=?, stock=? WHERE id=? RETURNING *")
				mock.ExpectQuery("UPDATE books SET title=?, slug=?, cover_image=?, synopsis=?, price=?, stock=? WHERE id=? RETURNING *").
					WithArgs(postUpdateBook.Title, postUpdateBook.Slug, postUpdateBook.CoverImage, postUpdateBook.Synopsis, postUpdateBook.Price, postUpdateBook.Stock, postUpdateBook.ID).
					WillReturnError(fmt.Errorf("error updating product"))

				_, err := r.Update(ctx, postUpdateBook)
				require.Error(t, err)

				// Ensure all expectations are met
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			withDatabaseMock(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				mockRepo := repositories.NewBookRepository(repositories.NewRepository(db))
				c.test(t, mockRepo, mock)
			})
		})
	}
}

func TestDeleteBook(t *testing.T) {
	cases := []struct {
		name string
		test func(*testing.T, *repositories.BookRepository, sqlmock.Sqlmock)
	}{
		{
			name: "Success",
			test: func(t *testing.T, r *repositories.BookRepository, mock sqlmock.Sqlmock) {
				ctx := context.Background()
				bookId := 1

				// Mock successful DELETE operation
				mock.ExpectExec("DELETE FROM books WHERE id=$1").WithArgs(bookId).WillReturnResult(
					sqlmock.NewResult(0, 1), // For DELETE, only rows affected (1) matters
				)

				err := r.Delete(ctx, int64(bookId))
				if err != nil {
					t.Fatalf("an error '%s' was not expected when deleting a book", err)
				}

				// Check that all expectations were met
				if err = mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			},
		},
		{
			name: "ExecContext Failure",
			test: func(t *testing.T, r *repositories.BookRepository, mock sqlmock.Sqlmock) {
				ctx := context.Background()
				bookId := 1

				// Mock failure during DELETE operation
				mock.ExpectExec("DELETE FROM books WHERE id=$1").WithArgs(bookId).WillReturnError(
					fmt.Errorf("error deleting book"),
				)

				err := r.Delete(ctx, int64(bookId))
				require.Error(t, err) // Assert that an error was returned

				// Check that all expectations were met
				if err = mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			},
		},
	}

	// Run all test cases
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			withDatabaseMock(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				mockRepo := repositories.NewBookRepository(repositories.NewRepository(db))
				c.test(t, mockRepo, mock)
			})
		})
	}
}
