/**
	server.go HTTP services for podkstr.com

	Routes:
		GET /	home

		Subs:
		GET /a/ admin
		GET /p/ private (user zone)
		GET /api/v1.0/ API


**/

package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
