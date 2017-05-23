package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/toorop/podkstr/core"
	"github.com/toorop/podkstr/logger"
)

// ResetPassword Reset Password
func ResetPassword(c echo.Context) error {
	type tpl struct {
		tplData
		Step    uint8
		Message string
		UUID    string
	}

	uuid := c.Param("uuid")

	var data tpl
	data.Title = "Activate your Podkstr account"
	data.Version = c.Get("version").(string)
	data.MoreScripts = []string{"vue.js", "axios.min.js", "components.js", "reset-password.js"}
	data.Step = 1

	if uuid != "" {
		// check if user exists
		_, found, err := core.UserGetByValidationUUID(uuid)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("%s - ResetPassword -> core.UserGetByValidationUUID(%s) - %s ", c.Request().RemoteAddr, uuid, err))
			return c.NoContent(http.StatusInternalServerError)
		}
		if !found {
			data.Step = 0
			data.Message = "This reset password link is unknow. Try reseting your password once again"
		} else {
			data.Step = 2
			data.UUID = uuid
		}
	}

	return c.Render(http.StatusOK, "resetpassword", data)
}
