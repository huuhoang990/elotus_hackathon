package main

import (
	"go_api/cmd/api/handlers"
)

func (app *Application) routes(handler handlers.Handler) {
	app.server.POST("/register", handler.RegisterHandler)
	app.server.POST("/login", handler.LoginHandler)
	app.server.GET("/auth", handler.GetAuthenticatedUser, app.appMiddleware.AuthMiddleware)
	app.server.POST("/upload", handler.UploadFile, app.appMiddleware.AuthMiddleware)
}
