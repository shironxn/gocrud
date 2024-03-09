package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Host string
		Port string
	}
	Database struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
	}
	JWTSecret string
}

var (
	cfg *Config
)

func NewConfig() (*Config, error) {
	if err := LoadConfig(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func LoadConfig() error {
	if cfg != nil {
		return nil
	}

	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("error loading environment variables: %w", err)
	}

	cfg = &Config{
		Server: struct {
			Host string
			Port string
		}{
			Host: os.Getenv("APP_HOST"),
			Port: os.Getenv("APP_PORT"),
		},
		Database: struct {
			Host     string
			Port     string
			Name     string
			User     string
			Password string
		}{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
		},
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	return nil
}
