package controllers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"github.com/toorop/podkstr/www/appContext"
)

// AjSignin login and sign up
func AjSignin(ec echo.Context) error {
	c := ec.(*appContext.AppContext)

	log.Info(viper.Get("apppath"))

	// Get a session
	session, err := c.GetCookieStore().Get(c.Request(), "podkastr")
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"ok": "false"})
	}

	// Set some session values.
	session.Values["foo"] = "bar"
	session.Values[42] = 43
	// Save it before we write to the response/return from the handler.
	session.Save(c.Request(), c.Response().Writer)
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}
