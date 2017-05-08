package controllers

import (
	"net/http"

	"github.com/labstack/echo"
)

// Todo controller
func Todo(c echo.Context) error {
	data := &tplData{
		Title:       "Podkstr",
		MoreScripts: []string{},
		UserEmail:   c.Get("uEmail").(string),
	}
	return c.Render(http.StatusOK, "todo", data)
}
