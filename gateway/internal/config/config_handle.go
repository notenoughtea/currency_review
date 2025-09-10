package config

import (
	"os"
	"strconv"

	"github.com/notenoughtea/currency_review/gateway/internal/logger"
	"gopkg.in/yaml.v3"
)

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type GRPC struct {
	Port     int    `yaml:"port"`
	Protocol string `yaml:"protocol"`
}

type AuthConfig struct {
	BaseURL string `yaml:"base_url"`
}

type Root struct {
	Server           Server     `yaml:"server"`
	CurrencyOuterKey string     `yaml:"CURRENCY_OUTER_KEY"`
	GRPC             GRPC       `yaml:"GRPC"`
	Auth             AuthConfig `yaml:"auth"`
}

var cfg Root

func Load() {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		logger.Log.Fatal("CONFIG_PATH not set")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		logger.Log.Fatal("read config:", err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		logger.Log.Fatal("unmarshal config:", err)
	}
	if p := os.Getenv("PORT"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			cfg.Server.Port = v
		}
	}
	if cfg.Server.Port == 0 {
		logger.Log.Fatal("server.port is 0")
	}
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}
}

func GetServerConfig() Server   { return cfg.Server }
func GetGrpcConfig() GRPC       { return cfg.GRPC }
func GetAuthConfig() AuthConfig { return cfg.Auth }
