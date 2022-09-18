package main

import (
	"home/controller"
	"home/storage"
	"net/http"

	"github.com/jasonlvhit/gocron"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	storage.OpenConnection()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", hello)

	gocron.Every(10).Seconds().Do(controller.CreateNews)
	<-gocron.Start()
	e.Logger.Fatal(e.Start(":8081"))
}
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
