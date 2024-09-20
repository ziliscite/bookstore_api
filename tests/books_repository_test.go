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

// Unit testing (or integration testing? Since its db?)
func TestCreateBook(t *testing.T) {
	book := &models.Book{
		Title: "Solo Leveling",
		// We insert the slug manually in this mocking
		// Because doing slug is handler's concern
		Slug:       "solo-leveling",
		CoverImage: "https://asura.nacmcdn.com/wp-content/uploads/2022/03/solo-leveling.jpg",
		Synopsis:   "In a world where hunters, humans with magical abilities, must battle deadly monsters to protect the human race, Sung Jinwoo, the weakest of the rank E hunters, struggles to earn a living. However, after narrowly surviving an overwhelmingly powerful dungeon that nearly wipes out his entire party, a mysterious program called the System chooses him as its sole player and grants him the extremely rare ability to level up in strength, possibly beyond any known limits. Jinwoo sets off on a journey to become an unparalleled S-rank hunter.",
		Price:      12.99,
		Stock:      517,
	}

	cases := []struct {
		name string
		test func(*testing.T, *repositories.Repository, sqlmock.Sqlmock)
	}{
		{
			name: "Success",
			test: func(t *testing.T, r *repositories.Repository, mock sqlmock.Sqlmock) {
				ctx := context.Background()
				bookId := 1

				mock.ExpectExec("INSERT INTO books (title, slug, cover_image, synopsis, price, stock) VALUES (?, ?, ?, ?, ?, ?)").WillReturnResult(
					sqlmock.NewResult(int64(bookId), 1),
				)

				createdBook, err := r.Create(ctx, book)
				if err != nil {
					t.Fatalf("an error '%s' was not expected when creating a book", err)
				}

				require.Equal(t, bookId, createdBook.ID)
				require.Equal(t, book.Title, createdBook.Title)
				require.Equal(t, book.Price, createdBook.Price)
				require.Equal(t, book.Stock, createdBook.Stock)

				if err = mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			},
		},
		{
			name: "Exec Failure",
			test: func(t *testing.T, r *repositories.Repository, mock sqlmock.Sqlmock) {
				ctx := context.Background()

				mock.ExpectExec("INSERT INTO books (title, slug, cover_image, synopsis, price, stock) VALUES (?, ?, ?, ?, ?, ?)").WillReturnError(
					fmt.Errorf("error while inserting book"),
				)

				_, err := r.Create(ctx, book)
				require.Error(t, err)

				if err = mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			},
		},
		{
			name: "Insert Id Failure",
			test: func(t *testing.T, r *repositories.Repository, mock sqlmock.Sqlmock) {
				ctx := context.Background()

				mock.ExpectExec("INSERT INTO books (title, slug, cover_image, synopsis, price, stock) VALUES (?, ?, ?, ?, ?, ?)").WillReturnError(
					fmt.Errorf("error getting last inserted id"),
				)

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
				mockRepo := repositories.NewRepository(db)
				c.test(t, mockRepo, mock)
			})
		})
	}
}

func TestGetBook(t *testing.T) {
	book := &models.Book{
		Title:      "Solo Leveling",
		Slug:       "solo-leveling",
		CoverImage: "https://asura.nacmcdn.com/wp-content/uploads/2022/03/solo-leveling.jpg",
		Synopsis:   "In a world where hunters, humans with magical abilities, must battle deadly monsters to protect the human race, Sung Jinwoo, the weakest of the rank E hunters, struggles to earn a living. However, after narrowly surviving an overwhelmingly powerful dungeon that nearly wipes out his entire party, a mysterious program called the System chooses him as its sole player and grants him the extremely rare ability to level up in strength, possibly beyond any known limits. Jinwoo sets off on a journey to become an unparalleled S-rank hunter.",
		Price:      12.99,
		Stock:      517,
	}

	cases := []struct {
		name string
		test func(*testing.T, *repositories.Repository, sqlmock.Sqlmock)
	}{
		{
			name: "Success",
			test: func(t *testing.T, r *repositories.Repository, mock sqlmock.Sqlmock) {
				ctx := context.Background()
				bookId := 1

				rows := sqlmock.NewRows([]string{
					"id", "title", "slug", "cover_image", "synopsis", "price", "stock",
				}).AddRow(bookId, book.Title, book.Slug, book.CoverImage, book.Synopsis, book.Price, book.Stock)

				mock.ExpectQuery("SELECT * FROM books WHERE id=?").WithArgs(bookId).WillReturnRows(rows)
				responseBook, err := r.GetById(ctx, bookId)
				if err != nil {
					t.Fatalf("an error '%s' was not expected when creating a book", err)
				}

				require.Equal(t, bookId, responseBook.ID)
				require.Equal(t, book.Title, responseBook.Title)
				require.Equal(t, book.Price, responseBook.Price)
				require.Equal(t, book.Stock, responseBook.Stock)

				if err = mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			},
		},
		{
			name: "GetContext Failure",
			test: func(t *testing.T, r *repositories.Repository, mock sqlmock.Sqlmock) {
				ctx := context.Background()

				mock.ExpectQuery("SELECT * FROM books WHERE id=?").WithArgs(1).WillReturnError(fmt.Errorf("error getting book"))
				_, err := r.GetById(ctx, 1)
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
				mockRepo := repositories.NewRepository(db)
				c.test(t, mockRepo, mock)
			})
		})
	}
}

func TestGetBooks(t *testing.T) {
	books := []models.Book{
		{
			Title:      "Omniscient Reader",
			Slug:       "omniscient-reader",
			CoverImage: "https://asura.nacmcdn.com/wp-content/uploads/2021/12/omniscient-reader.jpg",
			Synopsis:   "Dokja was an average office worker whose sole interest was reading his favorite web novel 'Three Ways to Survive the Apocalypse.' One day, his story comes to life, and he becomes the only person who knows how the world will end. He embarks on a journey to change the course of events and save humanity.",
			Price:      10.99,
			Stock:      432,
		},
		{
			Title:      "The Beginning After The End",
			Slug:       "the-beginning-after-the-end",
			CoverImage: "https://asura.nacmcdn.com/wp-content/uploads/2022/04/the-beginning-after-the-end.jpg",
			Synopsis:   "King Grey has unrivaled strength, wealth, and prestige in a world governed by martial ability. However, beneath the glamorous exterior lies a man devoid of purpose and will. Reincarnated into a new world filled with magic and monsters, Grey must now navigate this new world and uncover the secrets it holds.",
			Price:      11.49,
			Stock:      289,
		},
		{
			Title:      "Return of the Disaster-Class Hero",
			Slug:       "return-of-the-disaster-class-hero",
			CoverImage: "https://asura.nacmcdn.com/wp-content/uploads/2023/06/return-of-the-disaster-class-hero.jpg",
			Synopsis:   "In a world where powerful individuals known as 'heroes' stand at the top of society, Lee Geon was once regarded as the strongest disaster-class hero before being betrayed and left to die. After a long absence, he returns to seek revenge and reclaim his title.",
			Price:      13.79,
			Stock:      356,
		},
	}

	cases := []struct {
		name string
		test func(*testing.T, *repositories.Repository, sqlmock.Sqlmock)
	}{
		{
			name: "Success",
			test: func(t *testing.T, r *repositories.Repository, mock sqlmock.Sqlmock) {
				ctx := context.Background()

				rows := sqlmock.NewRows([]string{
					"id", "title", "slug", "cover_image", "synopsis", "price", "stock",
				})

				for i, book := range books {
					rows.AddRow(i+1, book.Title, book.Slug, book.CoverImage, book.Synopsis, book.Price, book.Stock)
				}

				mock.ExpectQuery("SELECT * FROM books").WillReturnRows(rows)
				responseBooks, err := r.GetAll(ctx)
				if err != nil {
					t.Fatalf("an error '%s' was not expected when creating a book", err)
				}

				for i, responseBook := range responseBooks {
					require.Equal(t, i+1, responseBook.ID)
					require.Equal(t, books[i].Title, responseBook.Title)
					require.Equal(t, books[i].Price, responseBook.Price)
					require.Equal(t, books[i].Stock, responseBook.Stock)
				}

				if err = mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			},
		},
		{
			name: "SelectContext Failure",
			test: func(t *testing.T, r *repositories.Repository, mock sqlmock.Sqlmock) {
				ctx := context.Background()

				mock.ExpectQuery("SELECT * FROM books").WillReturnError(fmt.Errorf("error getting books"))
				_, err := r.GetAll(ctx)
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
				mockRepo := repositories.NewRepository(db)
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
		Synopsis:   "In a world where hunters, humans with magical abilities, must battle deadly monsters to protect the human race, Sung Jinwoo, the weakest of the rank E hunters, struggles to earn a living. However, after narrowly surviving an overwhelmingly powerful dungeon that nearly wipes out his entire party, a mysterious program called the System chooses him as its sole player and grants him the extremely rare ability to level up in strength, possibly beyond any known limits. Jinwoo sets off on a journey to become an unparalleled S-rank hunter.",
		Price:      12.99,
		Stock:      517,
	}

	postUpdateBook := &models.Book{
		Title:      "Solo Leveling",
		Slug:       "solo-leveling",
		CoverImage: "https://asura.nacmcdn.com/wp-content/uploads/2022/03/solo-leveling.jpg",
		Synopsis:   "In a world where hunters, humans with magical abilities, must battle deadly monsters to protect the human race, Sung Jinwoo, the weakest of the rank E hunters, struggles to earn a living. However, after narrowly surviving an overwhelmingly powerful dungeon that nearly wipes out his entire party, a mysterious program called the System chooses him as its sole player and grants him the extremely rare ability to level up in strength, possibly beyond any known limits. Jinwoo sets off on a journey to become an unparalleled S-rank hunter.",
		Price:      15.99,
		Stock:      517,
	}

	cases := []struct {
		name string
		test func(*testing.T, *repositories.Repository, sqlmock.Sqlmock)
	}{
		{
			name: "Success",
			test: func(t *testing.T, r *repositories.Repository, mock sqlmock.Sqlmock) {
				ctx := context.Background()
				bookId := 1

				mock.ExpectExec("INSERT INTO books (title, slug, cover_image, synopsis, price, stock) VALUES (?, ?, ?, ?, ?, ?)").WillReturnResult(
					sqlmock.NewResult(int64(bookId), 1),
				)

				createdBook, err := r.Create(ctx, book)
				if err != nil {
					t.Fatalf("an error '%s' was not expected when creating a book", err)
				}

				require.Equal(t, createdBook.ID, bookId)
				require.Equal(t, createdBook, book)

				mock.ExpectExec("UPDATE books SET title=?, slug=?, cover_image=?, synopsis=?, price=?, stock=? WHERE id=?").WillReturnResult(
					sqlmock.NewResult(int64(bookId), 1),
				)

				updatedBook, err := r.Update(ctx, postUpdateBook)
				if err != nil {
					t.Fatalf("an error '%s' was not expected when updating a book", err)
				}

				require.Equal(t, updatedBook, postUpdateBook)
				if err = mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			},
		},
		{
			name: "NamedExecContext Failure",
			test: func(t *testing.T, r *repositories.Repository, mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE books SET title=?, slug=?, cover_image=?, synopsis=?, price=?, stock=? WHERE id=?").WillReturnError(
					fmt.Errorf("error updating product"),
				)

				_, err := r.Update(context.Background(), book)
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
				mockRepo := repositories.NewRepository(db)
				c.test(t, mockRepo, mock)
			})
		})
	}
}

func TestDeleteBook(t *testing.T) {
	cases := []struct {
		name string
		test func(*testing.T, *repositories.Repository, sqlmock.Sqlmock)
	}{
		{
			name: "Success",
			test: func(t *testing.T, r *repositories.Repository, mock sqlmock.Sqlmock) {
				ctx := context.Background()
				bookId := 1

				mock.ExpectExec("DELETE FROM books WHERE id=?").WithArgs(bookId).WillReturnResult(
					sqlmock.NewResult(int64(bookId), 1),
				)

				err := r.Delete(ctx, bookId)
				if err != nil {
					t.Fatalf("an error '%s' was not expected when creating a book", err)
				}

				if err = mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			},
		},
		{
			name: "ExecContext Failure",
			test: func(t *testing.T, r *repositories.Repository, mock sqlmock.Sqlmock) {
				ctx := context.Background()
				bookId := 1

				mock.ExpectExec("DELETE FROM books WHERE id=?").WithArgs(bookId).WillReturnError(
					fmt.Errorf("error deleting book"),
				)

				err := r.Delete(ctx, bookId)
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
				mockRepo := repositories.NewRepository(db)
				c.test(t, mockRepo, mock)
			})
		})
	}
}
