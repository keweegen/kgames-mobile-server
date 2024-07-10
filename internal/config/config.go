package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App        App
	Database   Database
	Redis      Redis
	GRPCServer GRPCServer
}

type App struct {
	Name        string `envconfig:"APP_NAME"`
	Version     string `envconfig:"APP_VERSION"`
	Environment string `envconfig:"APP_ENV"`
	Debug       bool   `envconfig:"APP_DEBUG"`
}

type Database struct {
	Host     string `envconfig:"DB_HOST"`
	Port     int    `envconfig:"DB_PORT"`
	Name     string `envconfig:"DB_NAME"`
	User     string `envconfig:"DB_USER"`
	Password string `envconfig:"DB_PASSWORD"`
	Schema   string `envconfig:"DB_SCHEMA"`
	SSLMode  string `envconfig:"DB_SSL_MODE"`
}

type Redis struct {
	Address string `envconfig:"REDIS_ADDRESS"`
	DB      int    `envconfig:"REDIS_DB"`
}

type GRPCServer struct {
	Port int `envconfig:"GRPC_SERVER_PORT"`
}

func (d Database) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.Schema, d.SSLMode)
}

func Read() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
