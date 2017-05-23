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
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"github.com/toorop/podkstr/core"
	"github.com/toorop/podkstr/logger"
	"github.com/toorop/podkstr/www/appContext"
	"github.com/toorop/podkstr/www/controllers"
)

// Version is podkstr webapp version
const Version = "0.1"

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
	// out
	logoutStr := viper.GetString("log.out")
	if logoutStr == "stdout" || logoutStr == "" {
		logger.Log.Out = os.Stdout
	} else {
		log.Println(logoutStr)
		f, err := os.OpenFile(logoutStr, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("unable to open log file ", logoutStr)
		}
		//defer f.Close()
		log.Println(f)
		logger.Log.Out = f
	}
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
	core.DB.LogMode(false)
	defer core.DB.Close()
	logger.Log.Info("database instantiated")

	/////////////////
	// Init Storer
	err = core.InitOsStore()
	if err != nil {
		logger.Log.Fatal("unable to init OsStore - ", err)
	}
	logger.Log.Info("openstack storer initialized ", core.Store)
	_, err = core.Store.Get("toto")
	logger.Log.Info("GET storer ", err)

	/////////////////
	// launch task runner
	taskRunner := core.NewTaskRunner()
	go taskRunner.Run()
	logger.Log.Info("taskruner Launched")

	//taskRunner.Stop()

	/////////////////
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
			cc.Set("version", Version)
			cc.SetCookieStore(sessions.NewCookieStore([]byte(viper.GetString("cookie.secret"))))
			return h(cc)
		}
	})

	/////////////////
	// Middlewares
	// checkuser
	e.Use(checkUser())
	// Logger
	loggerConfig := middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format: `${time_rfc3339_nano} - id: ${id} - remote_ip: ${remote_ip} - host: ${host} - ` +
			`method": ${method} - uri: ${uri} - status: ${status} - latency: ${latency} - ` +
			`latency_human: ${latency_human} - bytes_in: ${bytes_in} - ` +
			`bytes_out: ${bytes_out}` + "\n",
		Output: os.Stdout,
	}
	e.Use(middleware.LoggerWithConfig(loggerConfig))
	// recover
	e.Use(Recover())
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

	// Activate by validating email address
	e.GET("/activate/:uuid", controllers.ActivateAccount)

	// Reset Passwd
	e.GET("reset-password", controllers.ResetPassword)
	e.GET("reset-password/:uuid", controllers.ResetPassword)

	// Feed
	e.GET("/feed/:uuid", controllers.GetShowFeed)

	// private

	// dashboard
	e.GET("/dashboard", controllers.Dashboard)

	// AJAX

	// signin signup
	e.POST("/ajsignin", controllers.AjSignin)

	// send reset password email
	e.POST("/ajsendresetpasswordemail", controllers.AjSendResetPasswordEmail)

	// reset password
	e.POST("/ajresetpassword", controllers.AjResetPassword)

	// renvoi le mail d'activation
	e.POST("/ajresendactivationemail", controllers.AjResendActivationEmail)

	// Import Show
	e.POST("/ajimportshow", controllers.AjImportShow)

	// Delete show
	e.DELETE("/aj/show/delete/:uuid", controllers.AjDeleteShow)

	// Get User Shows
	e.GET("/aj/user/shows", controllers.AjGetUserShows)

	// Discourse SSO
	e.GET("/discourse/sso", controllers.DiscourseSSO)

	/////////////////
	// 10.9.8.7...0!
	logger.Log.Info("launch http")
	e.Logger.Fatal(e.Start(":1323"))
}
