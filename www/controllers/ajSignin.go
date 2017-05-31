package controllers

import (
	"bytes"
	"html/template"
	"net/http"
	"net/url"
	"path"

	"gopkg.in/gomail.v2"

	"golang.org/x/crypto/bcrypt"

	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"github.com/toorop/podkstr/core"
	"github.com/toorop/podkstr/logger"
	"github.com/toorop/podkstr/www/appContext"
)

// AjSignin login and sign up
func AjSignin(ec echo.Context) error {
	c := ec.(*appContext.AppContext)
	type FormData struct {
		Email   string `json:"email"`
		Passwd  string `json:"passwd"`
		Passwd2 string `json:"passwd2"`
		Signup  bool   `json:"signup"`
	}

	type response struct {
		Ok  bool
		Msg string
	}

	var err error
	var session *sessions.Session
	var resp = response{}
	var found bool
	var user core.User

	// Bind DATA
	fd := new(FormData)
	if err = c.Bind(&fd); err != nil {
		logger.Log.Error(c.Request().RemoteAddr, " - AjSignup -  ", err)
		return c.JSON(http.StatusOK, resp)
	}
	// validation
	fd.Email = strings.TrimSpace(strings.ToLower(fd.Email))
	if fd.Email == "" {
		logger.Log.Info(c.Request().RemoteAddr, " - Signup - bad request email is missing")
		return c.NoContent(http.StatusBadRequest)
	}

	fd.Passwd = strings.TrimSpace(strings.ToLower(fd.Passwd))
	if fd.Passwd == "" {
		logger.Log.Info(c.Request().RemoteAddr, " - Signup - bad request passwd is missing")
		return c.NoContent(http.StatusBadRequest)
	}

	// Login
	// Signup
	if fd.Signup {
		fd.Passwd2 = strings.TrimSpace(strings.ToLower(fd.Passwd2))
		if fd.Passwd2 == "" {
			logger.Log.Info(c.Request().RemoteAddr, " - Signup - bad request passwd2 is missing")
			return c.NoContent(http.StatusBadRequest)
		}
		_, found, err := core.UserGetByMail(fd.Email)
		if err != nil {
			logger.Log.Error(c.Request().RemoteAddr, " - Signup - models.GetUserByEmailPasswd - ", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		if found {
			resp.Msg = " There already and account associated with this email"
			return c.JSON(http.StatusOK, resp)
		}
		user, err = core.UserNew(fd.Email, fd.Passwd)
		if err != nil {
			logger.Log.Error(c.Request().RemoteAddr, " - Signup - models.UserNew - ", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		// Confirmation email
		mailTpl, err := template.ParseFiles(path.Join(viper.GetString("rootPath"), "etc/tpl/signup-confirm.eml"))
		if err != nil {
			logger.Log.Error(c.Request().RemoteAddr, " - Signup - sendmail template.ParseFiles() - ", err)
			return c.NoContent(http.StatusInternalServerError)
		}

		tplData := struct {
			Link string
		}{
			Link: viper.GetString("baseurl") + "/activate/" + url.QueryEscape(user.ValidationUUID),
		}

		buf := bytes.Buffer{}
		if err = mailTpl.Execute(&buf, tplData); err != nil {
			logger.Log.Error(c.Request().RemoteAddr, " - Signup - sendmail mailTpl.Execute - ", err)
			return c.NoContent(http.StatusInternalServerError)
		}

		m := gomail.NewMessage()
		m.SetAddressHeader("From", viper.GetString("smtp.sender"), "Podkstr")
		m.SetAddressHeader("To", user.Email, user.Email)
		m.SetHeader("Subject", "Welcome to Podkstr !")
		m.SetBody("text/plain", buf.String())

		d := gomail.NewPlainDialer(viper.GetString("smtp.host"), viper.GetInt("smtp.port"), viper.GetString("smtp.user"), viper.GetString("smtp.passwd"))
		if err = d.DialAndSend(m); err != nil {
			logger.Log.Error(c.Request().RemoteAddr, " - Signup - sendmail d.DialAndSend(m) - ", err)
		}

	} else {
		// sigin
		// Get user
		user, found, err = core.UserGetByEmailPasswd(fd.Email, fd.Passwd)
		if err != nil && err != bcrypt.ErrMismatchedHashAndPassword {
			logger.Log.Error(c.Request().RemoteAddr, " - Signup - models.GetUserByEmailPasswd - ", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		if !found || (err != nil && err == bcrypt.ErrMismatchedHashAndPassword) {
			logger.Log.Info(c.Request().RemoteAddr, " - Signup - user ", fd.Email, " not found or auth failed")
			resp.Msg = " Auth failed !"
			return c.JSON(http.StatusOK, resp)
		}
	}

	// Get a session
	session, err = c.GetCookieStore().Get(c.Request(), "podkstr")
	if err != nil {
		logger.Log.Error(c.Request().RemoteAddr, " - Signup - c.GetCookieStore().Get() - ", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	// Set some session values.
	session.Values["u@"] = user.Email

	// Save it before we write to the response/return from the handler.
	session.Save(c.Request(), c.Response().Writer)
	resp.Ok = true
	resp.Msg = "/dashboard"
	return c.JSON(http.StatusOK, resp)
}
