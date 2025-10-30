// Package config содержит типы конфигурации приложения и утилиты для их загрузки.
package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config агрегирует все разделы конфигурации приложения.
type Config struct {
	*AppConfig   `yaml:"app"`
	*RetryConfig `yaml:"retry"`
	*HTTPConfig  `yaml:"http"`
	*MongoConfig `yaml:"mongo"`
	*KafkaConfig `yaml:"kafka"`
}

// AppConfig описывает флаги включения подсистем приложения.
type AppConfig struct {
	IsHTTPEnabled   bool          `yaml:"http_enable" env:"APP_HTTP_ENABLED" env-default:"true"`
	IsGrpcEnabled   bool          `yaml:"grpc_enable" env:"APP_GRPC_ENABLED" env-default:"false"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"10s"`
}

// Mongo хранит параметры подключения к MongoDB.
type MongoConfig struct {
	Addr           string        `yaml:"addr"`
	Username       string        `yaml:"username"`
	Password       string        `yaml:"password"`
	DB             string        `yaml:"db_name"`
	ConnectTimeout time.Duration `yaml:"connect_timeout"`
	MaxPoolSize    uint64        `yaml:"max_pool_size"`
}

// KafkaConfig содержит настройки брокера Kafka, необходимые для инициализации
// продюсеров, консьюмеров и управления топиками.
type KafkaConfig struct {
	Address       string      `yaml:"address"`
	TestTopic     string      `yaml:"test-topic"`
	GroupID       string      `yaml:"group-id"`
	Network       string      `yaml:"network"`
	FetchBackoff  RetryConfig `yaml:"fetchBackoff"`
	CommitBackoff RetryConfig `yaml:"commitBackoff"`
}

// RetryConfig определяет параметры для механизма повторных попыток.
type RetryConfig struct {
	Attempts int           `yaml:"attempts" env-default:"3"`
	Initial  time.Duration `yaml:"initial"  env-default:"1s"`
	Max      time.Duration `yaml:"max"      env-default:"30s"`
	Factor   float64       `yaml:"factor"   env-default:"2.0"`
	Jitter   bool          `yaml:"jitter"   env-default:"true"`
}

// HTTPConfig задает настройки HTTP-сервера.
type HTTPConfig struct {
	Addr              string        `yaml:"addr" env:"HTTP_ADDR"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout" env:"HTTP_READ_HEADER_TIMEOUT"`
	ReadTimeout       time.Duration `yaml:"read_timeout" env:"HTTP_READ_TIMEOUT"`
	WriteTimeout      time.Duration `yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT"`
	IdleTimeout       time.Duration `yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT"`
}

// LoadConfig считывает конфигурацию из YAML-файла по указанному пути либо паникует при ошибке.
func LoadConfig(path string) *Config {
	stat, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			panic("config not found")
		}
		panic(fmt.Errorf("cannot access config: %w", err))
	}

	if stat.Size() == 0 {
		panic("config is empty")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}
