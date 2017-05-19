package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/toorop/podkstr/core"
	"github.com/toorop/podkstr/logger"
)

// ActivateAccount active account
func ActivateAccount(c echo.Context) error {
	/*if c.Get("userEmail") != nil {
		return c.Redirect(http.StatusPermanentRedirect, "/dashboard")
	}*/

	uuid := c.Param("uuid")
	logger.Log.Debug("UUID ", uuid)

	type tpl struct {
		tplData
		ValidationOk bool
	}

	var err error
	var data tpl
	data.Title = "Activate your Podkstr account"
	data.Version = c.Get("version").(string)
	data.MoreScripts = []string{"vue.js", "axios.min.js", "components.js", "activate-account.js"}
	data.ValidationOk = false

	// check uuid
	var user core.User
	user, data.ValidationOk, err = core.UserGetByValidationUUID(uuid)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("%s - ActivateAccount -> core.UserGetByValidationUUID(%s) - %s ", c.Request().RemoteAddr, uuid, err))
		return c.NoContent(http.StatusInternalServerError)
	}

	if data.ValidationOk {
		user.Activated = true
		if err = user.Save(); err != nil {
			logger.Log.Error(fmt.Sprintf("%s - ActivateAccount -> user.Save() - %s ", c.Request().RemoteAddr, err))
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	return c.Render(http.StatusOK, "activate", data)
}
