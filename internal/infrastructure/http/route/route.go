package route

import (
	"bookstore_api/internal/core/service"
	"bookstore_api/internal/infrastructure/http/controller"
	handlers "bookstore_api/internal/infrastructure/http/handler"
	"bookstore_api/internal/port"
	"bookstore_api/internal/repositories"
	"bookstore_api/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
)

type Router struct {
	Mux *chi.Mux
}

func NewRouter() *Router {
	// Initialize the router
	mux := chi.NewRouter()

	// Use middleware
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	return &Router{
		Mux: mux,
	}
}

func (r *Router) RegisterBookRoutes(repository port.BookRepository) {
	bookService, err := service.NewBookService(repository)
	if err != nil {
		log.Fatal(err)
	}

	bookHandler := controller.NewBookHandler(bookService)

	// Register book route
	r.Mux.Post("/books", bookHandler.CreateBook)
	r.Mux.Get("/books", bookHandler.GetAllBooks)
	r.Mux.Get("/books/{id}", bookHandler.GetBookById)
	r.Mux.Put("/books/{id}", bookHandler.UpdateBook)
	r.Mux.Delete("/books/{id}", bookHandler.DeleteBook)
}

func (r *Router) RegisterUserRoutes(handler *handlers.Handler, service *services.Service, repository *repositories.Repository) {
	userRepository := repositories.NewUserRepository(repository)
	userService := services.NewUserService(service, userRepository)

	sessionRepository := repositories.NewSessionRepository(repository)
	sessionService := services.NewSessionService(service, sessionRepository)

	userHandler := handlers.NewUserHandler(handler, userService, sessionService)

	r.Mux.Post("/register", userHandler.RegisterUser)
	r.Mux.Post("/login", userHandler.LoginUser)
	r.Mux.Post("/logout", userHandler.LogoutUser)

	r.Mux.Put("/update", userHandler.UpdateUser)

	r.Mux.Post("/refresh", userHandler.RefreshAccessToken)
	r.Mux.Put("/revoke", userHandler.RevokeAccessToken)
}
