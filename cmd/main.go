package main

import (
	"gocrud/internal/adapter/http/handler"
	"gocrud/internal/adapter/http/middleware"
	"gocrud/internal/adapter/http/route"
	"gocrud/internal/adapter/repository"
	"gocrud/internal/config"
	"gocrud/internal/core/domain"
	"gocrud/internal/core/service"
	"gocrud/internal/util"

	"github.com/charmbracelet/log"
	_ "github.com/shironxn/gocrud/docs"
)

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
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

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository, bcrypt)
	userHandler := handler.NewUserHandler(userService, validator, jwt)

	noteRepository := repository.NewNoteRepository(db)
	noteService := service.NewNoteService(noteRepository)
	noteHandler := handler.NewNoteHandler(noteService, validator)

	authMiddleware := middleware.NewAuthMiddleware(jwt)

	welcomeRoute := route.NewWelcomeRoute()
	authRoute := route.NewAuthRoute(userHandler, authMiddleware)
	userRoute := route.NewUserRoute(userHandler, authMiddleware)
	noteRoute := route.NewNoteRoute(noteHandler, authMiddleware)

	welcomeRoute.Route(app)
	authRoute.Route(app)
	userRoute.Route(app)
	noteRoute.Route(app)

	if err = app.Listen(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
		log.Fatal(err)
	}
}
