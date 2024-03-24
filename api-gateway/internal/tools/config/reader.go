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

type connections struct {
	Parsing string `yaml:"parsing"`
}

type gRPCConfig struct {
	Connections connections `yaml:"connections"`
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

// GetDomain получает домен сервера.
func (ac *appConfig) GetDomain() string { return ac.ServerConfig.Domain }

// GetPort получает порт сервера.
func (ac *appConfig) GetPort() string { return ac.ServerConfig.Port }

// GetTimeout получает время отклика сервера.
func (ac *appConfig) GetTimeout() time.Duration { return ac.ServerConfig.Timeout }

// GetLoggerLevel получает уровень логирования.
func (c *appConfig) GetLoggerLevel() string { return c.LoggerConfig.Level }

// GetGRPCParsingConnection - получает ссылку для подключения к сервису парсинга.
func (c *appConfig) GetGRPCParsingConnection() string { return c.GRPCConfig.Connections.Parsing }
