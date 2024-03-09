package main

import (
	"gocrud/internal/adapter/handler"
	"gocrud/internal/adapter/http/middleware"
	"gocrud/internal/adapter/http/route"
	"gocrud/internal/adapter/repository"
	"gocrud/internal/config"
	"gocrud/internal/core/domain"
	"gocrud/internal/core/service"
	"gocrud/internal/util"

	"github.com/charmbracelet/log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := config.NewGorm()
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&domain.User{})

	validate, err := util.NewValidator()
	if err != nil {
		log.Fatal(err)
	}

	bcrypt := util.NewBcrypt()
	jwt := util.NewJWT(cfg)

	app := config.NewFiber()

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository, bcrypt)
	userHandler := handler.NewUserHandler(userService, validate, jwt)

	authMiddleware := middleware.NewAuthMiddleware(jwt)

	authRoute := route.NewAuthRoute(userHandler)
	userRoute := route.NewUserRoute(userHandler, authMiddleware)

	authRoute.Route(app)
	userRoute.Route(app)

	if err = app.Listen(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
		log.Fatal(err)
	}
}
