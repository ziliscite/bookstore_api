package main

import (
	"bookstore_api/internal/handlers"
	"bookstore_api/internal/repositories"
	"bookstore_api/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	mux *chi.Mux
}

func NewRouter() *Router {
	// Initialize the router
	mux := chi.NewRouter()

	// Use middleware
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	return &Router{
		mux: mux,
	}
}

func (r *Router) RegisterBookRoutes(handler *handlers.Handler, repository *repositories.Repository) {
	// Setup book service and handler
	bookRepository := repositories.NewBookRepository(repository)
	bookService := services.NewBookService(bookRepository)
	bookHandler := handlers.NewBookHandler(handler, bookService)

	// Register book routes
	r.mux.Post("/books", bookHandler.CreateBook)
	r.mux.Get("/books", bookHandler.GetAllBooks)
	r.mux.Get("/books/{id}", bookHandler.GetBookById)
	r.mux.Put("/books/{id}", bookHandler.UpdateBook)
	r.mux.Delete("/books/{id}", bookHandler.DeleteBook)
}
