package controllers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/toorop/podkstr/core"
)

// Show display a show
func Show(c echo.Context) error {
	var userEmail string
	u := c.Get("user")
	if u != nil {
		userEmail = u.(core.User).Email
	}
	data := &tplData{
		Title:       "Podkstr",
		MoreScripts: []string{"vue.js", "axios.min.js", "components.js", "howler.js", "audioplayer.js", "show.js"},
		UserEmail:   userEmail,
		Version:     c.Get("version").(string),
	}
	return c.Render(http.StatusOK, "show", data)

}
