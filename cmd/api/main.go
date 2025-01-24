package main

import (
	"flag"
	"template/internal/repository"
	database "template/internal/repository/dbrepo"
	"time"

	"log"

	"github.com/gin-gonic/gin"
)

type application struct {
	Domain string
	DSN    string
	Host   string
	auth   Auth
	repo   repository.DatabaseRepository
}

func main() {
	// set application config
	var app application

	token_expiry := 0
	refresh_token_expiry := 0
	database_url := ""
	app.Domain = "example.com"
	// read from command line
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=movies sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.StringVar(&app.Host, "host", "localhost:8081", "host address")
	flag.StringVar(&database_url, "db_url", "mysql://root:123456789@127.0.0.1:3306/go_db", "data base url.sqlite://tmp/test.db or mysql://root:123456789@127.0.0.1:3306/go_db or postgres://admin:root@localhost:5432/test_db ")
	flag.StringVar(&app.auth.Secret, "jwt-secret", "verysecret", "signing secret")
	flag.StringVar(&app.auth.Issuer, "jwt-issuer", "example.com", "signing issuer")
	flag.StringVar(&app.auth.Audience, "jwt-audience", "example.com", "signing audience")
	flag.StringVar(&app.auth.CookieDomain, "cookie-domain", "localhost", "cookie domain")
	flag.StringVar(&app.auth.CookieDomain, "domain", "example.com", "domain")
	flag.IntVar(&token_expiry, "token_expiry", 15, "token expiry in minutes")
	flag.IntVar(&refresh_token_expiry, "refresh_token_expiry", 24, "refresh token expiry in hour")
	app.auth.CookiePath = "/"
	app.auth.CookieName = "__Host-refresh_token"
	app.auth.TokenExpiry = time.Minute * time.Duration(token_expiry)
	app.auth.RefreshExpiry = time.Hour * time.Duration(refresh_token_expiry)

	flag.Parse()

	// connect to the database
	app.repo = &database.GormDatabase{}
	err := app.repo.Init(database_url)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer app.repo.DeInit()

	// start a web server
	server := gin.Default()
	log.Println("server running on ", app.Host)
	app.routes(server)
	err = server.Run(app.Host)
	if err != nil {
		log.Fatal(err)
	}
}
