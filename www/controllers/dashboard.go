package controllers

import (
	"net/http"

	"github.com/labstack/echo"
)

// Dashboard controller
func Dashboard(c echo.Context) error {
	// chek auth
	if c.Get("uEmail").(string) == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/signin")
	}
	data := &tplData{
		Title:       "Podkstr dashboard",
		MoreScripts: []string{"vue.js", "axios.min.js", "components.js", "dashboard.js"},
		UserEmail:   c.Get("uEmail").(string),
	}
	return c.Render(http.StatusOK, "dashboard", data)
}
