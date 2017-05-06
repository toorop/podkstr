package controllers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

// SignIn login and sign up
func SignIn(c echo.Context) error {
	data := struct {
		Title       string
		MoreScripts []string
	}{
		Title:       "Sign in or Sign up",
		MoreScripts: []string{"vue.js", "axios.js", "components.js", "signin.js"},
	}
	log.Info(c.Request().Method)
	return c.Render(http.StatusOK, "signin", data)
}
