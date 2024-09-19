package main

import (
	"bookstore_api/db"
	"bookstore_api/internal/repositories"
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

	repository := &repositories.Repository{Db: database.Db}

	routers := NewRouter()
	routers.RegisterBookRoutes(repository)

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
