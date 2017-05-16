package controllers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/toorop/podkstr/core"
)

// Dashboard controller
func Dashboard(c echo.Context) error {
	// chek auth
	u := c.Get("user")
	if u == nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/signin")
	}
	data := &tplData{
		Title:       "Podkstr dashboard",
		MoreScripts: []string{"vue.js", "axios.min.js", "components.js", "dashboard.js"},
		UserEmail:   u.(core.User).Email,
		Version:     c.Get("version").(string),
	}
	return c.Render(http.StatusOK, "dashboard", data)
}
