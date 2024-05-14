package config

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Host string
		Port string
		Dev  string
		Web  string
	}
	Database struct {
		Host string
		Port string
		Name string
		User string
		Pass string
	}
	JWT struct {
		Access  string
		Refresh string
	}
}

var (
	config *Config
)

func NewConfig() (*Config, error) {
	if err := LoadConfig(); err != nil {
		return nil, err
	}
	return config, nil
}

func LoadConfig() error {
	if config != nil {
		return nil
	}

	if flag.Lookup("test.v") != nil {
		if err := godotenv.Load("../../.env"); err != nil {
			return err
		}
	} else {
		if err := godotenv.Load(); err != nil {
			return err
		}
	}

	config = &Config{
		Server: struct {
			Host string
			Port string
			Dev  string
			Web  string
		}{
			Host: os.Getenv("APP_HOST"),
			Port: os.Getenv("APP_PORT"),
			Dev:  os.Getenv("APP_DEV"),
			Web:  os.Getenv("APP_WEB"),
		},
		Database: struct {
			Host string
			Port string
			Name string
			User string
			Pass string
		}{
			Host: os.Getenv("DB_HOST"),
			Port: os.Getenv("DB_PORT"),
			Name: os.Getenv("DB_NAME"),
			User: os.Getenv("DB_USER"),
			Pass: os.Getenv("DB_PASS"),
		},
		JWT: struct {
			Access  string
			Refresh string
		}{
			Access:  os.Getenv("JWT_ACCESS_SECRET"),
			Refresh: os.Getenv("JWT_REFRESH_SECRET"),
		},
	}

	return nil
}
