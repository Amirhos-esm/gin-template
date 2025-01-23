package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string `json:"msg,omitempty"`
	Error   string  `json:"error_msg,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func (app *application) ErrorJson(c *gin.Context, err error, code int) {
	c.JSON(code, Response{
		Error: err.Error(),
	})
}
func (app *application) SendJson(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Data: data,
	})
}
func (app *application) hellowolrd(c *gin.Context) {

	c.JSON(http.StatusOK, Response{
		Message: "hello world",
	})
}

func (app *application) home(c *gin.Context) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "go movie app",
		Version: "1.0.0",
	}
	c.JSON(http.StatusOK, Response{
		Data: payload,
	})
}
func (app *application) movies(c *gin.Context) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "go movie app",
		Version: "1.0.0",
	}
	app.SendJson(c, payload)

}

func (app *application) authenticate(c *gin.Context) {
	// read json payload


	//validate user against database

	//check password


	// create a jwt user
	user := jwtUser{
		ID:        1,
		FirstName: "firstname",
		LastName:  "lastname",
	}
	tokens, err := app.auth.GenerateTokenPair(&user)
	if err != nil {
		app.ErrorJson(c, err, http.StatusInternalServerError)
		return
	}

	refreshCokie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(c.Writer, refreshCokie)

	app.SendJson(c, tokens)

}
