package main

import (
	"fmt"
	log "log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"parsing-service/internal/handlers"
	"parsing-service/internal/tools/config"
	parsinggrpc "parsing-service/internal/tools/grpc/parsing"
	"parsing-service/internal/tools/parsing"
	"parsing-service/internal/tools/postgres"
	"parsing-service/pkg/dbconn/pgsql"
	"parsing-service/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// Получение конфига.
	cfg, err := config.NewConfig("./configs/parsing-config.yaml")
	if err != nil {
		log.Error(fmt.Sprintf("ошибка прочтения файла конфигруаций: %v", err))
		return
	}

	// Настройка логов.
	logger.SetLogger(cfg.(config.LoggerConfig).LoggerLevel())

	// Создание парсера.
	parse := parsing.NewParsing(cfg.(config.ServerConfig).Timeout())

	// Соединение с БД.
	conn := pgsql.ConnectToDB(cfg.(config.DatabaseConfig).DSN())
	if conn == nil {
		log.Error("ошибка подключения к Postgres")
		return
	}

	db, err := postgres.NewDB(conn, cfg)
	if err != nil {
		log.Error(fmt.Sprintf("ошибка подключения к базе данных: %v", err))
		return
	}

	log.Info("Запуск parsing service")

	// Настройка конфигурации сервера.
	listenGRPC, err := net.Listen("tcp",
		fmt.Sprintf("%s:%s", cfg.(config.ServerConfig).Domain(), cfg.(config.ServerConfig).Port()))
	if err != nil {
		log.Error(fmt.Sprintf("ошибка прослушивания порта gRPC: %v", err))
		return
	}

	creds, err := credentials.NewServerTLSFromFile("./tls/cert.pem", "./tls/key.pem")
	if err != nil {
		log.Error(fmt.Sprintf("ошибка TLS подключения: %v", err))
		return
	}

	grpcSrv := grpc.NewServer(grpc.Creds(creds))
	parsinggrpc.RegisterParsingServer(grpcSrv, handlers.NewParser(db, parse))

	go func() {
		err = grpcSrv.Serve(listenGRPC)
		if err != nil {
			log.Error(fmt.Sprintf("ошибка запуска сервера gRPC: %v", err))
			os.Exit(1)
		}
	}()

	// Завершение работы.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	<-shutdown
	log.Info("Завершение работы...")

	grpcSrv.GracefulStop()
	err = listenGRPC.Close()
	if err != nil {
		log.Error(fmt.Sprintf("ошибка при завершении работы gRPC сервера %v", err))
	}

	err = conn.Close()
	if err != nil {
		log.Error(fmt.Sprintf("ошибка при завершении работы базы данных %v", err))
	}
}
