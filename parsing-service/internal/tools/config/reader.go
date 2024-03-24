package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"

	"github.com/caarlos0/env"
)

// appConfig - конфиг приложения.
type appConfig struct {
	ServerConfig     serverConfig     `yaml:"server"`
	DatabaseConfig   databaseConfig   `yaml:"database"`
	LoggerConfig     loggerConfig     `yaml:"logger"`
	EventsURIsConfig eventsURIsConfig `yaml:"events_uris"`
}

// serverConfig - конфиг сервера.
type serverConfig struct {
	Domain  string        `yaml:"domain"`
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// databaseConfig - конфиг базы данных.
type databaseConfig struct {
	DSN     string        `env:"DSN"`
	Timeout time.Duration `yaml:"timeout"`
}

// loggerConfig - структура конфиграции логов.
type loggerConfig struct {
	Level string `yaml:"level"`
}

type eventsURIsConfig struct {
	RussiaKrasnodar string `yaml:"russia_krasnodar"`
}

// NewConfig создает API конфига.
func NewConfig(configPath string) (Config, error) {
	// Считывание файла конфигурации.
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Декодирование данных в структуру конфигурации.
	var cfg appConfig
	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		return nil, err
	}

	err = env.Parse(&cfg.DatabaseConfig)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// GetDomain получает домен сервера.
func (ac *appConfig) GetDomain() string { return ac.ServerConfig.Domain }

// GetPort получает порт сервера.
func (ac *appConfig) GetPort() string { return ac.ServerConfig.Port }

// GetTimeout получает время отклика сервера.
func (ac *appConfig) GetTimeout() time.Duration { return ac.ServerConfig.Timeout }

// GetDSN получает строку для с данными для подключения к базе данных.
func (ac *appConfig) GetDSN() string { return ac.DatabaseConfig.DSN }

// GetDBTimeout получает время отклика базы данных.
func (ac *appConfig) GetDBTimeout() time.Duration { return ac.DatabaseConfig.Timeout }

// GetLoggerLevel получает уровень логирования.
func (c *appConfig) GetLoggerLevel() string { return c.LoggerConfig.Level }

// GetRussiaKrasnodar получает uri на события в месте.
func (c *appConfig) GetRussiaKrasnodar() string { return c.EventsURIsConfig.RussiaKrasnodar }
