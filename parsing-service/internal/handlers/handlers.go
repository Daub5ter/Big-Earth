package handlers

import (
	"net/http"
	"parsing-service/internal/data"
	"parsing-service/pkg/code"
)

type Parsing interface {
	Parse(r data.Place) *data.PlaceInformation
}

func Parse(p Parsing) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := data.Place{
			Country: "Russia",
			City:    "Krasnodar",
		}

		placeInformation := p.Parse(req)

		payload := code.JSONResponse{
			Error:   false,
			Message: "parsed",
			Data:    placeInformation,
		}

		code.WriteJSON(w, http.StatusOK, payload)
	}
}
