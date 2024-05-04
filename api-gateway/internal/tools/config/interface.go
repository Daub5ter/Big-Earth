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
	Domain() string
	Port() string
	Timeout() time.Duration
}

type GRPCConfig interface {
	GRPCParsingConnection() string
	GRPCParsingTimeout() time.Duration
}

// LoggerConfig - конфигурация логов приложения.
type LoggerConfig interface {
	LoggerLevel() string
}
