// Маршрутизатор.

package server

import (
	"net/http"
	"time"

	"api-gateway/internal/handlers"
	"api-gateway/internal/tools/config"
	"api-gateway/internal/tools/server/viewer"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// serv - структура сервера.
type serv struct {
	vr                viewer.Viewer
	parsingConnection string
	parsingTimeout    time.Duration
}

// NewServer создает новый сервер.
func NewServer(cfg config.GRPCConfig) Server {
	return serv{
		vr:                viewer.NewParamsViewer(),
		parsingConnection: cfg.GRPCParsingConnection(),
		parsingTimeout:    cfg.GRPCParsingTimeout(),
	}
}

// Routes возвращает обработчик с настройками и эндпоинтами.
func (s serv) Routes() http.Handler {
	// Создание и настройка маршрутизатора.
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(middleware.Heartbeat("/ping"))

	// Создание эндпоинтов.
	// тест времени и работы сервиса.
	//r.Get("/health", handlers.Health(s.database, s.pr))

	// grpcparsing
	r.Get("/parse/{country}/{city}", handlers.Parse(s.vr, s.parsingConnection, s.parsingTimeout))

	return r
}
