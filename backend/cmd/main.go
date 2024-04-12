package main

import (
	"github.com/shironxn/gocrud/internal/adapter/http/handler"
	"github.com/shironxn/gocrud/internal/adapter/http/middleware"
	"github.com/shironxn/gocrud/internal/adapter/http/route"
	"github.com/shironxn/gocrud/internal/adapter/repository"
	"github.com/shironxn/gocrud/internal/config"
	"github.com/shironxn/gocrud/internal/core/domain"
	"github.com/shironxn/gocrud/internal/core/service"
	"github.com/shironxn/gocrud/internal/util"

	_ "github.com/shironxn/gocrud/docs"

	"github.com/charmbracelet/log"
)

// @title gocrud
// @version 1.0
// @description golang crud api
// @host localhost:3000
// @BasePath /api/v1
func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := config.NewGorm(cfg).Connection()
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&domain.User{}, &domain.Note{})

	validator, err := util.NewValidator()
	if err != nil {
		log.Fatal(err)
	}

	app := config.NewFiber()
	bcrypt := util.NewBcrypt()
	jwt := util.NewJWT(cfg)
	pagination := util.NewPagination(*validator)

	userRepository := repository.NewUserRepository(db, pagination)
	userService := service.NewUserService(userRepository, bcrypt)
	userHandler := handler.NewUserHandler(userService, *validator, jwt)

	noteRepository := repository.NewNoteRepository(db, pagination)
	noteService := service.NewNoteService(noteRepository)
	noteHandler := handler.NewNoteHandler(noteService, *validator, jwt)

	authMiddleware := middleware.NewAuthMiddleware(jwt)

	initRoute := route.NewInitRoute()
	authRoute := route.NewAuthRoute(userHandler, authMiddleware)
	userRoute := route.NewUserRoute(userHandler, authMiddleware)
	noteRoute := route.NewNoteRoute(noteHandler, authMiddleware)

	initRoute.Route(app)
	authRoute.Route(app)
	userRoute.Route(app)
	noteRoute.Route(app)

	if err = app.Listen(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
		log.Fatal(err)
	}
}
