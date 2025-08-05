package main

import (
	"fmt"
	"github.com/labstack/echo/v4/middleware"
	"go_api/cmd/api/handlers"
	"go_api/cmd/api/middlewares"
	"go_api/cmd/common"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type Application struct {
	logger        echo.Logger
	server        *echo.Echo
	handler       handlers.Handler
	appMiddleware middlewares.AppMiddleware
}

func main() {
	e := echo.New()
	err := godotenv.Load()

	if err != nil {
		e.Logger.Fatal("Error loading .env file")
	}

	db, err := common.NewMysql()
	if err != nil {
		e.Logger.Fatal(err)
	}

	h := handlers.Handler{
		DB:     db,
		Logger: e.Logger,
	}
	appMiddleware := middlewares.AppMiddleware{
		DB:     db,
		Logger: e.Logger,
	}
	app := Application{
		logger:        e.Logger,
		server:        e,
		handler:       h,
		appMiddleware: appMiddleware,
	}
	e.Use(middleware.Logger(), middlewares.CustomMiddleware)
	app.routes(h)

	port := os.Getenv("APP_PORT")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
