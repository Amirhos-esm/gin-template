package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Auth struct {
	Issuer        string
	Audience      string
	Secret        string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
	CookieDomain  string
	CookiePath    string
	CookieName    string
}

type jwtUser struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
type TokenPairs struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type Claims struct {
	jwt.RegisteredClaims
}

func (j *Auth) GenerateTokenPair(user *jwtUser) (*TokenPairs, error) {
	// Create a token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = fmt.Sprintf("%s %s",user.FirstName,user.LastName)
	claims["sub"] = fmt.Sprint(user.ID)
	claims["aud"] = j.Audience
	claims["iss"] = j.Issuer
	claims["iat"] = time.Now().UTC().Unix()
	claims["typ"] = "JWT"

	// Set the expiry for jwt
	claims["exp"] = time.Now().UTC().Add(j.TokenExpiry).Unix()

	// Create a signed token
	signedAccessToken , err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return nil,err
	}
	
	// Set the refresh claims
	token = jwt.New(jwt.SigningMethodHS256)
	claims = token.Claims.(jwt.MapClaims)
	claims["sub"] = fmt.Sprint(user.ID)
	claims["iat"] = time.Now().UTC().Unix()


	// Set the expiry for jwt
	claims["exp"] = time.Now().UTC().Add(j.RefreshExpiry).Unix()

	// Create a signed token
	signedRefreshToken , err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return nil,err
	}

	return &TokenPairs{
		Token: signedAccessToken,
		RefreshToken: signedRefreshToken,
	},nil

}


func (j *Auth) GetRefreshCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name: j.CookieName,
		Path: j.CookiePath,
		Value: refreshToken,
		Expires: time.Now().Add(j.RefreshExpiry),
		MaxAge: int(j.RefreshExpiry.Seconds()),
		SameSite: http.SameSiteStrictMode,
		Domain: j.CookieDomain,
		HttpOnly: true,
		Secure: true,
	}
}

func (j *Auth) GetExpiredRefreshCookie() *http.Cookie {
	return &http.Cookie{
		Name: j.CookieName,
		Path: j.CookiePath,
		Value: "",
		Expires: time.Unix(0, 0),
		MaxAge: -1,
		SameSite: http.SameSiteStrictMode,
		Domain: j.CookieDomain,
		HttpOnly: true,
		Secure: true,
	}
}