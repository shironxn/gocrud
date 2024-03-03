package config

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

type App struct {
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

var app *App

func GetConfig() *App {
	if app == nil {
		initConfig()
	}

	return app
}

func initConfig() *App {
	config := App{}

	err := godotenv.Load(".env")
	if err != nil {
		log.Error(err)
	}

	config.Server.Host = os.Getenv("APP_HOST")
	config.Server.Port = os.Getenv("APP_PORT")

	config.Database.Host = os.Getenv("DB_HOST")
	config.Database.Port = os.Getenv("DB_PORT")
	config.Database.Name = os.Getenv("DB_NAME")
	config.Database.User = os.Getenv("DB_USER")
	config.Database.Pass = os.Getenv("DB_PASS")

	config.JWTSecret = os.Getenv("JWT_SECRET")

	return &config
}
