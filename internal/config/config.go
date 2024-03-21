package config

import (
	"errors"
	"flag"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Host string
		Port string
	}
	Database struct {
		Host string
		Port string
		Name string
		User string
		Pass string
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

func LoadEnv(fileName string) error {
	err := godotenv.Load(fileName)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		if err := godotenv.Load(".env.example"); err != nil {
			return errors.New("failed to load environment file")
		}

		return err
	}

	return nil
}

func LoadConfig() error {
	if cfg != nil {
		return nil
	}

	if flag.Lookup("test.v") != nil {
		if err := godotenv.Load("../../.env.test"); err != nil {
			return err
		}
	} else {
		if err := LoadEnv(".env"); err != nil {
			return err
		}
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
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	return nil
}
