package controllers

import (
	"net/http"

	"github.com/labstack/echo"
)

// Home / controller
func Home(c echo.Context) error {
	return c.Render(http.StatusOK, "home", "titi")
}
