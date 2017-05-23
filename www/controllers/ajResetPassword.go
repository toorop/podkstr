package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/toorop/podkstr/core"
	"github.com/toorop/podkstr/logger"
	"github.com/toorop/podkstr/www/appContext"
)

// AjResetPassword reset user password
func AjResetPassword(ec echo.Context) error {
	c := ec.(*appContext.AppContext)
	var err error
	type response struct {
		Ok  bool
		Msg string
	}
	resp := new(response)

	type FormData struct {
		UUID   string `json:"uuid"`
		Passwd string `json:"passwd"`
	}

	fd := new(FormData)
	if err = c.Bind(&fd); err != nil {
		logger.Log.Error(fmt.Sprintf("%s - AjResetPassword -> c.Bind(&fd) - %s ", c.Request().RemoteAddr, err))
		return c.JSON(http.StatusInternalServerError, resp)
	}

	if len(fd.Passwd) < 6 {
		resp.Msg = "your password must be at least 6 chars lenght"
		return c.JSON(http.StatusOK, resp)
	}

	// get user from validation UUID
	user, found, err := core.UserGetByValidationUUID(fd.UUID)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("%s - AjResetPassword -> core.UserGetByValidationUUID(%s) - %s ", c.Request().RemoteAddr, fd.UUID, err))
		return c.JSON(http.StatusInternalServerError, resp)
	}
	if !found {
		resp.Msg = "no such user"
		return c.JSON(http.StatusOK, resp)
	}

	if err = user.SetPasswd(fd.Passwd); err != nil {
		logger.Log.Error(fmt.Sprintf("%s - AjResetPassword -> user.SetPassword(%s) - %s ", c.Request().RemoteAddr, fd.Passwd, err))
		return c.JSON(http.StatusInternalServerError, resp)
	}

	// reset activation uuid
	if err = user.ResetValidationUUID(); err != nil {
		logger.Log.Error(fmt.Sprintf("%s - AjResetPassword -> user.ResetValidationUUID() - %s ", c.Request().RemoteAddr, err))
		return c.JSON(http.StatusInternalServerError, resp)
	}
	resp.Ok = true
	return c.JSON(http.StatusOK, resp)
}
