package main

import (
	"context"
	"fmt"
	log "log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"api-gateway/internal/tools/config"
	"api-gateway/internal/tools/server"
	"api-gateway/pkg/logger"
)

func main() {
	// Получение конфига.
	cfg, err := config.NewConfig("./app/api-gateway-config.yaml")
	if err != nil {
		log.Error(fmt.Sprintf("ошибка прочтения файла конфигруаций: %v", err))
		return
	}

	// Настройка логов.
	logger.SetLogger(cfg.(config.LoggerConfig).GetLoggerLevel())

	log.Info("Запуск grpcparsing service")

	// Настройка конфигурации сервера.
	s := server.NewServer(cfg)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.(config.ServerConfig).GetDomain(), cfg.(config.ServerConfig).GetPort()),
		Handler:      s.Routes(),
		ReadTimeout:  cfg.(config.ServerConfig).GetTimeout(),
		WriteTimeout: cfg.(config.ServerConfig).GetTimeout(),
	}

	// Запуск сервера.
	go func() {
		err = srv.ListenAndServeTLS("./app/cert.pem", "./app/key.pem")
		if err != nil {
			log.Error(fmt.Sprintf("ошибка запуска сервера: %v", err))
			return
		}
	}()

	// Завершение работы.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	<-shutdown
	log.Info("Завершение работы...")

	err = srv.Shutdown(context.Background())
	if err != nil {
		log.Error(fmt.Sprintf("ошибка при завершении работы сервера %v", err))
	}
}
