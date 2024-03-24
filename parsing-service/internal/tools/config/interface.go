package config

import "time"

// Config - API для работы с конфигом.
type Config interface {
	ServerConfig
	DatabaseConfig
	LoggerConfig
	EventsURIs
}

// ServerConfig - API конфига сервера.
type ServerConfig interface {
	GetDomain() string
	GetPort() string
	GetTimeout() time.Duration
}

// DatabaseConfig - API конфига базы данных.
type DatabaseConfig interface {
	GetDSN() string
	GetDBTimeout() time.Duration
}

// LoggerConfig - конфигурация логов приложения.
type LoggerConfig interface {
	GetLoggerLevel() string
}

// EventsURIs - конфигурация uris событий места.
type EventsURIs interface {
	GetRussiaKrasnodar() string
}
