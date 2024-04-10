package config

import "time"

// Config - API для работы с конфигом.
type Config interface {
	ServerConfig
	DatabaseConfig
	LoggerConfig
}

// ServerConfig - API конфига сервера.
type ServerConfig interface {
	Domain() string
	Port() string
	Timeout() time.Duration
}

// DatabaseConfig - API конфига базы данных.
type DatabaseConfig interface {
	DSN() string
	DBTimeout() time.Duration
}

// LoggerConfig - конфигурация логов приложения.
type LoggerConfig interface {
	LoggerLevel() string
}
