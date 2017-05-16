package controllers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/toorop/podkstr/core"
)

// Home / controller
func Home(c echo.Context) error {
	var userEmail string
	u := c.Get("user")
	if u != nil {
		userEmail = u.(core.User).Email
	}
	data := tplData{
		Title:       "Podkstr",
		MoreScripts: []string{},
		UserEmail:   userEmail,
		Version:     c.Get("version").(string),
	}
	return c.Render(http.StatusOK, "home", data)
}
