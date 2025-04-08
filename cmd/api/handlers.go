package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Response struct {
	Message string `json:"msg,omitempty"`
	Error   string `json:"error_msg,omitempty"`
	Data    any    `json:"data,omitempty"`
}


func (app *application) hellowolrd(c *gin.Context) {
	id , err := GetPathParam[uuid.UUID](c,"id",func(u uuid.UUID) error {
		return nil
	})
	if err != nil {
		SendError(c,err,http.StatusBadRequest)
		return
	}
	
	c.JSON(http.StatusOK, Response{
		Message: fmt.Sprintf("hello world to %v",id),
	})
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
// authenticate handles user authentication by validating the provided email and password,
// generating a JWT token pair upon successful authentication, and setting a refresh token cookie.
//
// @Summary Authenticate user
// @Description Authenticate user with email and password, and return JWT token pair
// @Tags authentication
// @Accept json
// @Produce json
// @Param Authentication body Authentication true "User credentials"
// @Success 200 {object} TokenPairs "JWT token pair"
// @Router /authenticate [post]
func (app *application) authenticate(c *gin.Context) {
	// read json payload
	input := Authentication{}
	// Parse and validate the JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		SendError(c, err, http.StatusBadRequest)
		return
	}
	//validate user against database
	user, err := app.repo.GetUserByEmail(input.Email)
	if err != nil {
		SendError(c, err, http.StatusInternalServerError)
		return
	}
	if user == nil {
		SendError(c, nil, http.StatusUnauthorized)
		return
	}
	//check password
	if ok, err := user.PasswordMatches(input.Password); err != nil || !ok {
		SendError(c, nil, http.StatusUnauthorized)
		return
	}

	// create a jwt user
	jwt_user := jwtUser{
		ID:        1,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
	tokens, err := app.auth.GenerateTokenPair(&jwt_user)
	if err != nil {
		SendError(c, err, http.StatusInternalServerError)
		return
	}

	refreshCokie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(c.Writer, refreshCokie)

	SendJson(c, tokens)

}

func (app *application) refreshToken(c *gin.Context) {
	for _, cookie := range c.Request.Cookies() {
		if cookie.Name == app.auth.CookieName {
			claims := &Claims{}
			refreshToken := cookie.Value

			// parse the token to get the claims
			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.auth.Secret), nil
			})
			if err != nil {
				SendError(c, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			// get the user id from the token claims
			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				SendError(c, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			user, err := app.repo.GetUserById(uint64(userID))
			if err != nil {
				SendError(c, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}
			if user == nil {
				SendError(c, nil, http.StatusUnauthorized)
				return
			}

			u := jwtUser{
				ID:        int(user.ID),
				FirstName: user.FirstName,
				LastName:  user.LastName,
			}

			tokenPairs, err := app.auth.GenerateTokenPair(&u)
			if err != nil {
				SendError(c, errors.New("error generating tokens"), http.StatusUnauthorized)
				return
			}

			http.SetCookie(c.Writer, app.auth.GetRefreshCookie(tokenPairs.RefreshToken))

			SendJson(c, tokenPairs)
			return

		}
	}
	SendError(c, nil, http.StatusUnauthorized)
}

func (app *application) logout(c *gin.Context) {
	http.SetCookie(c.Writer, app.auth.GetExpiredRefreshCookie())
	c.Writer.WriteHeader(http.StatusAccepted)
}
