package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"

	"github.com/caarlos0/env"
)

// appConfig - конфиг приложения.
type appConfig struct {
	ServerConfig   serverConfig   `yaml:"server"`
	DatabaseConfig databaseConfig `yaml:"database"`
	LoggerConfig   loggerConfig   `yaml:"logger"`
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

// Domain получает домен сервера.
func (ac *appConfig) Domain() string { return ac.ServerConfig.Domain }

// Port получает порт сервера.
func (ac *appConfig) Port() string { return ac.ServerConfig.Port }

// Timeout получает время отклика сервера.
func (ac *appConfig) Timeout() time.Duration { return ac.ServerConfig.Timeout }

// DSN получает строку для с данными для подключения к базе данных.
func (ac *appConfig) DSN() string { return ac.DatabaseConfig.DSN }

// DBTimeout получает время отклика базы данных.
func (ac *appConfig) DBTimeout() time.Duration { return ac.DatabaseConfig.Timeout }

// LoggerLevel получает уровень логирования.
func (c *appConfig) LoggerLevel() string { return c.LoggerConfig.Level }
