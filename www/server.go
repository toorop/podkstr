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
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"github.com/toorop/podkstr/core"
	"github.com/toorop/podkstr/www/appContext"
	"github.com/toorop/podkstr/www/controllers"
	"github.com/toorop/podkstr/www/logger"
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

	viper.Set("rootPath", rootPath)
	// Load config
	viper.AddConfigPath(rootPath + "/etc")
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("unable to read config - ", err)
	}
	//log.Println("config loaded")

	// init app logger (! access log)
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = time.RFC3339Nano
	customFormatter.FullTimestamp = true
	logger.Log = logrus.New()
	logger.Log.Formatter = customFormatter
	logger.Log.Out = os.Stdout
	logger.Log.Level = logrus.DebugLevel
	logger.Log.Info("logrus instantiated")

	// Init DB
	core.DB, err = gorm.Open(viper.GetString("db.dialect"), viper.GetString("db.args"))
	if err != nil {
		logger.Log.Fatal("database connexion failed - ", err)
	}
	if err = core.DbAutoMigrate(); err != nil {
		logger.Log.Fatal("unable to automigrate DB - ", err)
	}
	core.DB.Debug()
	logger.Log.Info("database instantiated")

	// init echo web server
	e := echo.New()

	// Debug
	e.Debug = true

	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "header:X-XSRF-TOKEN",
		CookieName:  "XSRF-TOKEN",
	}))

	// Custom context
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := appContext.NewAppContext(c)
			cc.SetCookieStore(sessions.NewCookieStore([]byte(viper.GetString("cookie.secret"))))
			return h(cc)
		}
	})

	/////////////////
	// Middlewares
	// checkuser
	e.Use(checkUser())

	// Logger
	e.Use(middleware.Logger())
	// recover
	e.Use(middleware.Recover())
	// Log Error
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				logger.Log.Error(err)
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

	// Signout
	e.GET("/signout", controllers.Signout)

	// private

	// dashboard
	e.GET("/dashboard", controllers.Dashboard)

	// AJAX

	// signin signup
	e.POST("/ajsignin", controllers.AjSignin)

	// Import Show
	e.POST("/ajimportshow", controllers.AjImportShow)

	/////////////////
	// 10.9.8.7...0!
	e.Logger.Fatal(e.Start(":1323"))
}
