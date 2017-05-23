package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"

	"gopkg.in/gomail.v2"

	"strings"

	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"github.com/toorop/podkstr/core"
	"github.com/toorop/podkstr/logger"
	"github.com/toorop/podkstr/www/appContext"
)

// AjSendResetPasswordEmail send reset password email
func AjSendResetPasswordEmail(ec echo.Context) error {
	c := ec.(*appContext.AppContext)
	var err error
	type response struct {
		Ok  bool
		Msg string
	}
	resp := new(response)

	type FormData struct {
		Email string `json:"email"`
	}

	fd := new(FormData)
	if err = c.Bind(&fd); err != nil {
		logger.Log.Error(fmt.Sprintf("%s - AjSendResetPasswordEmail -> c.Bind(&fd) - %s ", c.Request().RemoteAddr, err))
		return c.JSON(http.StatusInternalServerError, resp)
	}
	user, found, err := core.UserGetByMail(strings.TrimSpace(fd.Email))
	if err != nil {
		logger.Log.Error(fmt.Sprintf("%s - AjSendResetPasswordEmail -> core.UserGetByMail(%s) - %s ", c.Request().RemoteAddr, fd.Email, err))
		return c.JSON(http.StatusInternalServerError, resp)
	}
	if !found {
		resp.Msg = fmt.Sprintf("no such user %s", fd.Email)
		return c.JSON(http.StatusOK, resp)
	}

	// generate new validation UUID
	if err = user.ResetValidationUUID(); err != nil {
		logger.Log.Error(fmt.Sprintf("%s - AjSendResetPasswordEmail -> user.ResetValidationUUID( - %s ", c.Request().RemoteAddr, err))
		return c.NoContent(http.StatusInternalServerError)
	}

	// Confirmation email
	mailTpl, err := template.ParseFiles(path.Join(viper.GetString("rootPath"), "etc/tpl/reset-password.eml"))
	if err != nil {
		logger.Log.Error(fmt.Sprintf("%s - AjSendResetPasswordEmail -> template.ParseFile - %s ", c.Request().RemoteAddr, err))
		return c.NoContent(http.StatusInternalServerError)
	}

	tplData := struct {
		Link string
	}{
		Link: viper.GetString("baseurl") + "/reset-password/" + url.QueryEscape(user.ValidationUUID),
	}

	buf := bytes.Buffer{}
	if err = mailTpl.Execute(&buf, tplData); err != nil {
		logger.Log.Error(fmt.Sprintf("%s - AjSendResetPasswordEmail -> mailTpl.Execute - %s ", c.Request().RemoteAddr, err))
		return c.NoContent(http.StatusInternalServerError)
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", viper.GetString("smtp.sender"), "Podkstr")
	m.SetAddressHeader("To", user.Email, user.Email)
	m.SetHeader("Subject", "Reset your Podkstr password")
	m.SetBody("text/plain", buf.String())

	d := gomail.NewPlainDialer(viper.GetString("smtp.host"), viper.GetInt("smtp.port"), viper.GetString("smtp.user"), viper.GetString("smtp.passwd"))
	if err = d.DialAndSend(m); err != nil {
		logger.Log.Error(fmt.Sprintf("%s - AjSendResetPasswordEmail -> d.DialAndSend(m) - %s ", c.Request().RemoteAddr, err))
		return c.NoContent(http.StatusInternalServerError)
	}

	resp.Ok = true
	resp.Msg = fmt.Sprintf("a email with our activation link has been sent to %s, check your mailbox", user.Email)
	return c.JSON(http.StatusOK, resp)
}
