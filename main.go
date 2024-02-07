package main

import (
	"log"
	"net/http"
	"os"
	henrymedscodechallenge "rob/henry-meds-code-challenge/src"
)

func main() {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	database := henrymedscodechallenge.CreateDB()

	router := henrymedscodechallenge.GetRouter(database, logger)

	logger.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", router)
}
