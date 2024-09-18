package main

import (
	"bookstore_api/db"
	"context"
	"github.com/joho/godotenv"
	"log"
)

func runServer() {
	context.TODO()
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database, err := db.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	defer func(database *db.Database) {
		err = database.Close()
		if err != nil {

		}
	}(database)
}
