package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Root struct {
	Server           Server `yaml:"server"`
	CurrencyOuterKey string `yaml:"CURRENCY_OUTER_KEY"`
	DB               DB     `yaml:"DB"`
	GRPC             GRPC   `yaml:"GRPC"`
}

type Server struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type GRPC struct {
	Port     int    `yaml:"port"`
	Protocol string `yaml:"protocol"`
}

type DB struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Port     int    `yaml:"port"`
	SSLMode  string `yaml:"sslmode"`
	TimeZone string `yaml:"TimeZone"`
}

func GetServerConfig() Server {
	data, err := os.ReadFile("../../../config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var root Root
	err = yaml.Unmarshal(data, &root)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return root.Server
}

func GetDbConfig() DB {
	data, err := os.ReadFile("../../../config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var root Root
	err = yaml.Unmarshal(data, &root)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return root.DB
}

func GetToken() string {
	data, err := os.ReadFile("../../../config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var root Root
	err = yaml.Unmarshal(data, &root)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return root.CurrencyOuterKey
}

func GetConfGRPC() GRPC {
	data, err := os.ReadFile("../../../config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var root Root
	err = yaml.Unmarshal(data, &root)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return root.GRPC
}
