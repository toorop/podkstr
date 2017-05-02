package controllers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

// Home / controller
func Home(c echo.Context) error {
	err := c.Render(http.StatusOK, "home", "titi")
	if err != nil {
		log.Error("ERROR ", err)

	}
	return err
}
