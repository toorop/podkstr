package controllers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/toorop/podkstr/www/appContext"
)

// AjPrivateTodo todo controller for ajax request in private zone
func AjPrivateTodo(ec echo.Context) error {
	type response struct {
		Ok  bool
		Msg string
	}
	var resp = response{}
	c := ec.(*appContext.AppContext)
	u := c.Get("user")
	if u == nil {
		resp.Msg = "You are not logged please signin"
		return c.JSON(http.StatusForbidden, resp)
	}
	resp.Ok = true
	resp.Msg = "TODO"
	return c.JSON(http.StatusForbidden, resp)
}
