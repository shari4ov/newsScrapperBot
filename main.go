package main

import (
	"home/controller"
	"home/storage"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	storage.OpenConnection()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", hello)

	e.GET("/create/news", controller.CreateNews)
	e.Logger.Fatal(e.Start(":8081"))
}
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
