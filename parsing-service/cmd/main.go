package main

import (
	"parsing-service/internal/db"
	"parsing-service/internal/handlers"
	"parsing-service/internal/parse"
)

func main() {
	parsing := parse.NewParsing()

	r := db.Request{
		Country: "Russia",
		City:    "Krasnodar",
	}

	handlers.Parse(parsing, r)
}
