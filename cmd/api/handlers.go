package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string `json:"msg,omitempty"`
	Error   string `json:"error_msg,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func (app *application) ErrorJson(c *gin.Context, err error, code int) {
	if err == nil {
		c.JSON(code, Response{
		})
		return
	}
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
	type Authentication struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	input := Authentication{}
	// Parse and validate the JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		app.ErrorJson(c, err, http.StatusBadRequest)
		return
	}
	//validate user against database
	user, err := app.repo.GetUserByEmail(input.Email)
	if err != nil {
		app.ErrorJson(c, err, http.StatusInternalServerError)
		return
	}
	if user == nil {
		app.ErrorJson(c, nil, http.StatusUnauthorized)
		return
	}
	//check password
	if ok, err := user.PasswordMatches(input.Password); err!=nil || !ok {
		app.ErrorJson(c, nil, http.StatusUnauthorized)
		return
	}

	// create a jwt user
	jwt_user := jwtUser{
		ID:        1,
		FirstName: user.Name,
		LastName:  user.LastName,
	}
	tokens, err := app.auth.GenerateTokenPair(&jwt_user)
	if err != nil {
		app.ErrorJson(c, err, http.StatusInternalServerError)
		return
	}

	refreshCokie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(c.Writer, refreshCokie)

	app.SendJson(c, tokens)

}
