package controllers

import (
	"net/http"

	"github.com/labstack/echo"
)

// SignIn login and sign up
func SignIn(c echo.Context) error {
	if c.Get("userEmail") != nil {
		return c.Redirect(http.StatusPermanentRedirect, "/dashboard")
	}
	data := tplData{
		Title:       "Sign in or Sign up",
		MoreScripts: []string{"vue.js", "axios.min.js", "components.js", "signin.js"},
		Version:     c.Get("version").(string),
	}
	return c.Render(http.StatusOK, "signin", data)
}
