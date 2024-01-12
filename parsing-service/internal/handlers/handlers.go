package handlers

import (
	"net/http"
	"parsing-service/internal/models"
	"parsing-service/pkg/code"
)

type ParsingI interface {
	Parse(models.Place) (*models.PlaceInformation, error)
}

func Parse(p ParsingI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := models.Place{
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
