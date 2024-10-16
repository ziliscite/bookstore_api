package main

import (
	"bookstore_api/internal/infrastructure/http/handler"
	"bookstore_api/internal/infrastructure/http/route"
	"bookstore_api/internal/infrastructure/postgres"
	"bookstore_api/internal/infrastructure/redis"
	"bookstore_api/internal/repositories"
	"bookstore_api/internal/services"
	"context"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func runServer() error {
	database, err := postgres.New()
	if err != nil {
		return err
	}

	defer func(database *postgres.Database) error {
		err = database.Close()
		if err != nil {
			return err
		}
		return nil
	}(database)

	cache := redis.New()
	defer func(cache *redis.Cache) error {
		err = cache.Close()
		if err != nil {
			return err
		}
		return nil
	}(cache)

	_, err = cache.GetCache().Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	handler := handler.NewHandler(cache.GetCache())
	service, err := services.NewService()
	if err != nil {
		return err
	}
	repository := repositories.NewRepository(database.GetDB())

	routers := route.NewRouter()
	routers.Mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("ok"))
		if err != nil {
			return
		}
	})

	// Diff
	bookRepo := postgres.NewBookRepository(database)
	routers.RegisterBookRoutes(bookRepo)
	//

	routers.RegisterUserRoutes(handler, service, repository)

	err = http.ListenAndServe(":8081", routers.Mux)
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
		log.Fatalf("Error starting server, %s", err)
	}
}

// docker exec -it bookstore-postgresql psql -U postgres
