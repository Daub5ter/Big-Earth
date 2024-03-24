package config

import "time"

// Config - API для работы с конфигом.
type Config interface {
	ServerConfig
	GRPCConfig
	LoggerConfig
}

// ServerConfig - API конфига сервера.
type ServerConfig interface {
	GetDomain() string
	GetPort() string
	GetTimeout() time.Duration
}

type GRPCConfig interface {
	GetGRPCParsingConnection() string
}

// LoggerConfig - конфигурация логов приложения.
type LoggerConfig interface {
	GetLoggerLevel() string
}
