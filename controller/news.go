package controller

import (
	"home/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateNews(c echo.Context) error {
	service.CreateNewsService()
	return c.String(http.StatusCreated, "create")

}
