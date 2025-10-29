package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	*AppConfig   `yaml:"app"`
	*RetryConfig `yaml:"retry"`
	*HTTPConfig  `yaml:"http"`
}

type AppConfig struct {
	IsHttpEnabled bool `yaml:"http_enable" env:"APP_HTTP_ENABLED" env-default:"true"`
	IsGrpcEnabled bool `yaml:"grpc_enable" env:"APP_GRPC_ENABLED" env-default:"false"`
}

// RetryConfig определяет параметры для механизма повторных попыток.
type RetryConfig struct {
	Attempts int           `yaml:"attempts" default:"3"`
	Initial  time.Duration `yaml:"initial"  default:"1s"`
	Max      time.Duration `yaml:"max"      default:"30s"`
	Factor   float64       `yaml:"factor"   default:"2.0"`
	Jitter   bool          `yaml:"jitter"   default:"true"`
}

type HTTPConfig struct {
	Addr              string        `yaml:"addr" env:"HTTP_ADDR"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout" env:"HTTP_READ_HEADER_TIMEOUT"`
	ReadTimeout       time.Duration `yaml:"read_timeout" env:"HTTP_READ_TIMEOUT"`
	WriteTimeout      time.Duration `yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT"`
	IdleTimeout       time.Duration `yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT"`
}

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
