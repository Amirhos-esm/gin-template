package main

import "github.com/gin-gonic/gin"

func (app *application) routes(server *gin.Engine) {
	server.Use(app.CORSMiddleware())

	server.GET("/", app.home)
	server.GET("/movies",app.movies)

	server.POST("/authenticate",app.authenticate)
	server.GET("/refresh",app.refreshToken)
	server.GET("/logout",app.logout)


	server.GET("/hello", app.AuthRequiredMiddleware() ,app.hellowolrd)

}
 