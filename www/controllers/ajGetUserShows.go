package controllers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/toorop/podkstr/core"
	"github.com/toorop/podkstr/logger"
	"github.com/toorop/podkstr/www/appContext"
)

// AjGetUserShows returns User Shows
func AjGetUserShows(ec echo.Context) error {
	type response struct {
		Ok    bool
		Msg   string
		Shows []core.Show
	}
	var resp = response{}
	var err error

	c := ec.(*appContext.AppContext)
	u := c.Get("user")
	if u == nil {
		resp.Msg = "You are not logged please signin"
		return c.JSON(http.StatusForbidden, resp)
	}

	// get User show
	shows, err := u.(core.User).GetShows()
	if err != nil {
		logger.Log.Error(c.Request().RemoteAddr, " - AjGetUserShows -  ", err)
		return c.JSON(http.StatusInternalServerError, resp)
	}
	resp.Ok = true
	resp.Shows = shows
	return c.JSON(http.StatusOK, resp)
}
