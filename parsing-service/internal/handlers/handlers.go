package handlers

import (
	"net/http"
	"parsing-service/internal/data"
	"parsing-service/pkg/code"
)

type Parsing interface {
	Parse(data.Place) (*data.PlaceInformation, error)
}

func Parse(p Parsing) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := data.Place{
			Country: "Russia",
			City:    "Krasnodar",
		}

		placeInformation, err := p.Parse(req)
		if err != nil {
			payload := code.JSONResponse{
				Error:   true,
				Message: "not parsed",
				Data:    err,
			}

			code.WriteJSON(w, http.StatusBadRequest, payload)
		}

		payload := code.JSONResponse{
			Error:   false,
			Message: "parsed",
			Data:    placeInformation,
		}

		code.WriteJSON(w, http.StatusOK, payload)
	}
}
