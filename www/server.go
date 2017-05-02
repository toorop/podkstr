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
	"html/template"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/toorop/podkstr/core"
	"github.com/toorop/podkstr/www/controllers"
)

func main() {
	var err error
	// Init DB
	core.DB, err = gorm.Open("sqlite3", "/home/toorop/Projects/Go/src/github.com/toorop/podkstr/etc/podkstr.db")
	if err != nil {
		log.Fatal("DB connexion failed -", err)
	}
	if err = core.DbAutoMigrate(); err != nil {
		log.Fatal("unable to automigrate DB", err)
	}

	// init echo web server
	e := echo.New()

	e.Debug = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//e.Use(middleware.Static(""))

	// Template
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = t

	// Routes

	e.Static("/static", "/home/toorop/Projects/Go/src/github.com/toorop/podkstr/www/dist/static")

	// Home
	e.GET("/home", controllers.Home)

	// Login / Signup

	e.Logger.Fatal(e.Start(":1323"))
}
