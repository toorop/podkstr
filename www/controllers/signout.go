package controllers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/toorop/podkstr/www/appContext"
	"github.com/toorop/podkstr/www/logger"
)

// Signout controller
func Signout(ec echo.Context) error {
	c := ec.(*appContext.AppContext)
	// Get a session
	session, err := c.GetCookieStore().Get(c.Request(), "podkastr")
	if err != nil {
		logger.Log.Error(c.Request().RemoteAddr, " - Signup - c.GetCookieStore().Get() - ", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	// Remove Sessions variables
	delete(session.Values, "u@")
	session.Save(c.Request(), c.Response().Writer)
	c.Set("uEmail", nil)
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
