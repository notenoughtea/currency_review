package config

import (
	"os"
	"strconv"

	"github.com/notenoughtea/currency_review/currency/internal/logger"
	"gopkg.in/yaml.v3"
)

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DB struct {
	Host            string `yaml:"host"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	DBName          string `yaml:"dbname"`
	Port            int    `yaml:"port"`
	SSLMode         string `yaml:"sslmode"`
	TimeZone        string `yaml:"TimeZone"`
	Connect_timeout int    `yaml:"Connect_timeout"`
}

type GRPC struct {
	Port     int    `yaml:"port"`
	Protocol string `yaml:"protocol"`
}

type Root struct {
	Server           Server `yaml:"server"`
	CurrencyOuterKey string `yaml:"CURRENCY_OUTER_KEY"`
	DB               DB     `yaml:"DB"`
	GRPC             GRPC   `yaml:"GRPC"`
}

var cfg Root

func Load() {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		logger.Log.Fatal("CONFIG_PATH not set")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		logger.Log.Fatalf("read config: %v", err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		logger.Log.Fatalf("unmarshal config: %v", err)
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

func GetServerConfig() Server { return cfg.Server }
func GetDbConfig() DB         { return cfg.DB }
func GetGrpcConfig() GRPC     { return cfg.GRPC }
func GetToken() string        { return cfg.CurrencyOuterKey }
