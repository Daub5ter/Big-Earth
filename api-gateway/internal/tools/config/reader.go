package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

// appConfig - конфиг приложения.
type appConfig struct {
	ServerConfig serverConfig `yaml:"server"`
	LoggerConfig loggerConfig `yaml:"logger"`
	GRPCConfig   gRPCConfig   `yaml:"grpc"`
}

// serverConfig - конфиг сервера.
type serverConfig struct {
	Domain  string        `yaml:"domain"`
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type parsing struct {
	Connection string        `yaml:"connection"`
	Timeout    time.Duration `yaml:"timeout"`
}

type gRPCConfig struct {
	Parsing parsing `yaml:"parsing"`
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

	return &cfg, nil
}

// Domain получает домен сервера.
func (ac *appConfig) Domain() string { return ac.ServerConfig.Domain }

// Port получает порт сервера.
func (ac *appConfig) Port() string { return ac.ServerConfig.Port }

// Timeout получает время отклика сервера.
func (ac *appConfig) Timeout() time.Duration { return ac.ServerConfig.Timeout }

// LoggerLevel получает уровень логирования.
func (c *appConfig) LoggerLevel() string { return c.LoggerConfig.Level }

// GRPCParsingConnection - получает ссылку для подключения к сервису парсинга.
func (c *appConfig) GRPCParsingConnection() string { return c.GRPCConfig.Parsing.Connection }

// GRPCParsingTimeout - получает таймаут для отклика сервиса парсинга.
func (c *appConfig) GRPCParsingTimeout() time.Duration { return c.GRPCConfig.Parsing.Timeout }
