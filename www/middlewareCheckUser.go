package main

import (
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/toorop/podkstr/core"
	"github.com/toorop/podkstr/logger"
	"github.com/toorop/podkstr/www/appContext"
)

// TODO HERE

func checkUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ec echo.Context) error {
			var err error
			var session *sessions.Session
			c := ec.(*appContext.AppContext)
			// Get a session
			session, err = c.GetCookieStore().Get(c.Request(), "podkastr")
			if err != nil {
				logger.Log.Error(c.Request().RemoteAddr, " - middlewareCheckuser - c.GetCookieStore().Get() - ", err)
				return c.NoContent(http.StatusInternalServerError)
			}
			email := session.Values["u@"]
			c.Set("uEmail", "")
			if email != nil && email != "" && !strings.HasPrefix(c.Request().RequestURI, "/static/") {
				user, found, err := core.UserGetByMail(email.(string))
				if err != nil {
					logger.Log.Error(c.Request().RemoteAddr, " - middlewareCheckuser - core.UserGetByMail for ", email, " - ", err)
					return c.NoContent(http.StatusInternalServerError)
				}
				if found {
					c.Set("user", user)
				} else {
					c.Set("uEmail", "")
					//logger.Log.Info(c.Request().RemoteAddr, " - middlewareCheckuser - core.UserGetByMail bad email in cookie: ", email)
					//return c.NoContent(http.StatusInternalServerError)
				}
			}
			return next(c)
		}
	}
}
