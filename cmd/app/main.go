package main

import (
	"bookstore_api/db"
	"bookstore_api/internal/handlers"
	"bookstore_api/internal/repositories"
	"bookstore_api/internal/services"
	"context"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func runServer() error {
	database, err := db.NewDatabase()
	if err != nil {
		return err
	}

	defer func(database *db.Database) error {
		err = database.Close()
		if err != nil {
			return err
		}
		return nil
	}(database)

	cache := db.NewCache()
	defer func(cache *db.Cache) error {
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

	handler := handlers.NewHandler(cache.GetCache())
	service, err := services.NewService()
	if err != nil {
		return err
	}
	repository := repositories.NewRepository(database.GetDB())

	routers := NewRouter()
	routers.mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("ok"))
		if err != nil {
			return
		}
	})

	routers.RegisterBookRoutes(handler, service, repository)

	err = http.ListenAndServe(":8080", routers.mux)
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
