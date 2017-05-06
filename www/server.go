/**
	server.go HTTP services for podkstr.com

	Routes:

		Subs:
		GET /a/ admin
		GET /p/ private (user zone)
		GET /api/v1.0/ API


**/

package main

import (
	"html/template"

	"os"
	"path/filepath"

	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"github.com/toorop/podkstr/core"
	"github.com/toorop/podkstr/www/appContext"
	"github.com/toorop/podkstr/www/controllers"
)

// Main
func main() {
	var err error
	var rootPath string

	// get root
	rootPath, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal("unable to get root path -", err)
	}

	// Load config
	viper.AddConfigPath(rootPath + "/etc")
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("unable to read config - ", err)
	}

	log.Info("config loaded")

	// Init DB
	core.DB, err = gorm.Open(viper.GetString("db.dialect"), viper.GetString("db.args"))
	if err != nil {
		log.Fatal("database connexion failed - ", err)
	}
	if err = core.DbAutoMigrate(); err != nil {
		log.Fatal("unable to automigrate DB - ", err)
	}
	log.Info("database instantiated")
	// init echo web server
	e := echo.New()

	// Debug
	e.Debug = true

	// Custom context
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := appContext.NewAppContext(c)
			cc.SetCookieStore(sessions.NewCookieStore([]byte("bal")))
			return h(cc)
		}
	})

	/////////////////
	// Middlewares

	// Logger
	e.Use(middleware.Logger())
	// recover
	e.Use(middleware.Recover())
	// Log Error
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				log.Error(err)
			}
			return err
		}
	})

	/////////////////
	// Templates

	t := &Template{
		templates: template.Must(template.ParseGlob(rootPath + "/views/*.html")),
	}
	e.Renderer = t

	/////////////////
	// Routes

	// Static
	e.Static("/static", rootPath+"/static")

	// Home
	e.GET("/", controllers.Home)

	// Signin / Sign up
	e.GET("/signin", controllers.SignIn)

	// Login / Signup

	// AJAX

	// signin signup
	e.POST("/ajsignin", controllers.AjSignin)

	/////////////////
	// 10.9.8.7...0!
	e.Logger.Fatal(e.Start(":1323"))
}
