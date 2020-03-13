package main

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/lmnzr/simpleshop/cmd/simpleshop/database"

	_ "github.com/lmnzr/simpleshop/cmd/simpleshop/docs"
	"github.com/lmnzr/simpleshop/cmd/simpleshop/hello"
	"github.com/lmnzr/simpleshop/cmd/simpleshop/helper/env"
	"github.com/lmnzr/simpleshop/cmd/simpleshop/helper/jwt"

	logutil "github.com/lmnzr/simpleshop/cmd/simpleshop/helper/log"
	"github.com/lmnzr/simpleshop/cmd/simpleshop/middleware"
	"github.com/lmnzr/simpleshop/cmd/simpleshop/types"
	log "github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logutil.Logger(nil).Error("Error loading .env file")
	}

	environment := env.Getenv("ENVIRONMENT", "development")

	if environment != "development" {
		log.SetLevel(log.InfoLevel)

		log.SetFormatter(&log.JSONFormatter{})

		log.SetOutput(&lumberjack.Logger{
			Filename:   "var/log/app.log",
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		})
	} else {
		log.SetLevel(log.DebugLevel)
	}

}

// @title Simpleshop Swagger API
// @version 1.0
// @description Swagger API for Golang Project Simpleshop.
// @termsOfService http://swagger.io/terms/
// @BasePath
func main() {

	router := echo.New()
	middleware.Setup(router)

	db, _ := database.OpenDbConnection()
	errping := db.Ping()

	if errping != nil {
		logutil.LoggerDB().Panic("failed to connect to database")
	}

	//Use New Custom Context With DB
	router.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &types.DBContext{c, db}
			return h(cc)
		}
	})

	port := env.GetenvI("PORT", 9000)

	router.GET("/", public)
	router.GET("/swagger/*", echoSwagger.WrapHandler)
	router.GET("/forbidden/", forbidden)
	router.GET("/protected/", protected, middleware.JwtMiddleware)
	router.GET("/credential/", credential)

	hello.Routes(router)

	lock := make(chan error)
	go func(lock chan error) { lock <- router.Start(fmt.Sprintf(":%d", port)) }(lock)

	err := <-lock
	if err != nil {
		logutil.Logger(nil).Panic("failed to start application")
	}
}

func public(c echo.Context) error {
	cc := c.(*types.DBContext)
	return cc.Context.String(http.StatusOK, "Welcome to simpleshop")
}

func forbidden(c echo.Context) error {
	return echo.NewHTTPError(500, "Forbidden Land")
}

func protected(c echo.Context) error {
	cc := c.(*types.DBContext)
	name := cc.Context.Get("name").(string)
	return c.String(http.StatusOK, "Welcome "+name)
}

func credential(c echo.Context) error {
	cc := c.(*types.DBContext)

	cred := jwt.Credential{
		Name:  "Almas",
		UUID:  "11037",
		Admin: true,
	}
	pl := jwt.NewPayload(cred)
	token, _ := jwt.Signing(pl)
	a := types.Auth{
		Token: token,
	}
	return cc.Context.JSON(http.StatusOK, a)
}