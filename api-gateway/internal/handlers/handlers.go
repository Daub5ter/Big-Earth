package handlers

import (
	"context"
	"errors"
	"fmt"
	log "log/slog"
	"net/http"
	"time"

	"api-gateway/internal/tools/grpc/parsing"
	"api-gateway/internal/tools/server/viewer"
	"api-gateway/pkg/code"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// errBadRequest - это ошибка, которая обрабатываются в функциях.
var (
	errBadRequest = errors.New("плохо сформулирован запрос")
	errCreds      = errors.New("ошибка сертефикатов")
	errConnection = errors.New("ошибка соединения с сервером")
)

// Parse получает данные из сервиса парсинга.
func Parse(vr viewer.Viewer, connection string, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		country := vr.ViewParam(r, "country")
		if country == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errBadRequest)
			return
		}

		city := vr.ViewParam(r, "city")
		if city == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errBadRequest)
			return
		}

		creds, err := credentials.NewClientTLSFromFile("./tls/cert.pem", "")
		if err != nil {
			log.Error(fmt.Sprintf("не создается клиент tls: %v", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errCreds)
			return
		}

		conn, err := grpc.Dial(connection, grpc.WithTransportCredentials(creds))
		if err != nil {
			log.Error(fmt.Sprintf("не создается клиент grpc: %v", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errCreds)
			return
		}
		defer func() { _ = conn.Close() }()

		c := parsing.NewParsingClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		placeInformation, err := c.Parse(ctx,
			&parsing.Place{
				Country: country,
				City:    city,
			},
		)
		if err != nil {
			log.Error(fmt.Sprintf("не отправляется запрос: %v", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errConnection)
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
