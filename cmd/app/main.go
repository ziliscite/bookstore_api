package main

import (
	"bookstore_api/db"
	"bookstore_api/internal/controllers"
	"bookstore_api/internal/repositories"
	"bookstore_api/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func runServer() error {
	database, err := db.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	defer func(database *db.Database) {
		err = database.Close()
		if err != nil {

		}
	}(database)

	bookRepo := &repositories.Repository{Db: database.Db}
	bookService := services.NewBookService(bookRepo)
	bookHandler := controllers.NewBookHandler(bookService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/books", bookHandler.CreateBook)
	r.Get("/books", bookHandler.GetAllBooks)
	r.Get("/books/{id}", bookHandler.GetBookById)
	r.Put("/books/{id}", bookHandler.UpdateBook)
	r.Delete("/books/{id}", bookHandler.DeleteBook)

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = runServer()
	if err != nil {
		log.Fatal("Error starting server")
	}
}
