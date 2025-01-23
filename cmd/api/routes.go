package main

import "github.com/gin-gonic/gin"

func (app *application) routes(server *gin.Engine) {
	server.Use(app.CORSMiddleware())

	server.GET("/", app.home)
	server.GET("/hello", app.hellowolrd)
	server.GET("/movies",app.movies)
	server.GET("/authenticate",app.authenticate)
}
 