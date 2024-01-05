package handlers

import (
	"net/http"
	"parsing-service/internal/data"
)

type Parsing interface {
	Parse(r data.Request)
}

func Parse(p Parsing) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := data.Request{
			Country: "Russia",
			City:    "Krasnodar",
		}
		p.Parse(req)
	}
}
