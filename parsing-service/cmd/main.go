package main

import (
	"fmt"
	"log"
	"net/http"
	"parsing-service/api"
	"parsing-service/internal/data"
	"time"
)

const webPort = "1234"
const timeOfWaiting = 15 * time.Second

func main() {
	parsing := data.NewParsing()

	log.Println("Starting parsing service")

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", webPort),
		Handler:      api.Routes(parsing),
		ReadTimeout:  timeOfWaiting,
		WriteTimeout: timeOfWaiting,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

	//r := data.Request{
	//	Country: "Russia",
	//	City:    "Krasnodar",
	//}
	//	handlers.Parse(parsing, r)
}
