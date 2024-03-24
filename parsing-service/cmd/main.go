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
	"parsing-service/internal/tools/database"
	"parsing-service/internal/tools/parsing"
	"parsing-service/internal/tools/rpc/grpcparsing"
	"parsing-service/pkg/dbconn/pgsql"
	"parsing-service/pkg/logger"

	"google.golang.org/grpc"
)

func main() {
	// Получение конфига.
	cfg, err := config.NewConfig("./app/parsing-config.yaml")
	if err != nil {
		log.Error(fmt.Sprintf("ошибка прочтения файла конфигруаций: %v", err))
		return
	}

	// Настройка логов.
	logger.SetLogger(cfg.(config.LoggerConfig).GetLoggerLevel())

	// Создание парсера.
	parse := parsing.NewParsing(cfg)

	// Соединение с БД.
	conn := pgsql.ConnectToDB(cfg.GetDSN())
	if conn == nil {
		log.Error("ошибка подключения к Postgres")
		return
	}

	db, err := database.NewDB(conn, cfg)
	if err != nil {
		log.Error(fmt.Sprintf("ошибка подключения к базе данных: %v", err))
		return
	}

	log.Info("Запуск parsing service")

	// Настройка конфигурации сервера.
	listenGRPC, err := net.Listen("tcp",
		fmt.Sprintf("%s:%s", cfg.(config.ServerConfig).GetDomain(), cfg.(config.ServerConfig).GetPort()))
	if err != nil {
		log.Error(fmt.Sprintf("ошибка прослушивания порта gRPC: %v", err))
		return
	}

	grpcSrv := grpc.NewServer()
	grpcparsing.RegisterParsingServer(grpcSrv, handlers.NewParser(db, parse))

	go func() {
		err = grpcSrv.Serve(listenGRPC)
		if err != nil {
			log.Error(fmt.Sprintf("ошибка запуска сервера gRPC: %v", err))
			return
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

	err = pgsql.CloseConnection(conn)
	if err != nil {
		log.Error(fmt.Sprintf("ошибка при завершении работы базы данных %v", err))
	}
}
