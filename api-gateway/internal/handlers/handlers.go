package handlers

import (
	"context"
	"net/http"

	"api-gateway/internal/models"
	"api-gateway/internal/tools/rpc/grpcparsing"
	"api-gateway/internal/tools/server/viewer"
	"api-gateway/pkg/code"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Parse получает данные из сервиса парсинга.
func Parse(vr viewer.Viewer, connection string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		country := vr.ViewParam(r, "country")
		if country == "" {
			code.ErrorJSON(w, http.StatusBadRequest, models.ErrBadRequest)
			return
		}

		city := vr.ViewParam(r, "city")
		if city == "" {
			code.ErrorJSON(w, http.StatusBadRequest, models.ErrBadRequest)
			return
		}

		conn, err := grpc.Dial(connection, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		defer func() { _ = conn.Close() }()

		c := grpcparsing.NewParsingClient(conn)

		// todo: изменить context на нормальную реализацию с таймаутом или отменой.
		placeInformation, err := c.Parse(context.Background(),
			&grpcparsing.Place{
				Country: country,
				City:    city,
			},
		)
		if err != nil {
			code.ErrorJSON(w, http.StatusInternalServerError, err)
			return
		}

		payload := code.JSONResponse{
			Error:   false,
			Message: "parsed",
			Data:    placeInformation,
		}

		code.WriteJSON(w, http.StatusOK, payload)
	}
}
