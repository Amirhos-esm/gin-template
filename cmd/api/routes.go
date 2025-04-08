package main

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files" // gin-swagger middleware
	"github.com/swaggo/gin-swagger"
	_ "template/docs"
)

func (app *application) routes(server *gin.Engine) {
	if app.ExposeOpenApi {
		server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	server.Use(app.CORSMiddleware())

	server.POST("/authenticate", app.authenticate)
	server.GET("/refresh", app.refreshToken)
	server.GET("/logout", app.logout)

	server.GET("/hello/:id", app.AuthRequiredMiddleware(), app.hellowolrd)
	server.GET("/id/:id", app.hellowolrd)


}
