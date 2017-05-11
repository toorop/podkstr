package controllers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/toorop/podkstr/core"
	"github.com/toorop/podkstr/www/appContext"
	"github.com/toorop/podkstr/www/logger"
)

// AjDeleteShow delete User Show
// TODO ne pas supprimer si il est en cours de synch ->  mettre un lock
func AjDeleteShow(ec echo.Context) error {
	type response struct {
		Ok  bool
		Msg string
	}
	var resp = response{}
	var err error

	c := ec.(*appContext.AppContext)
	u := c.Get("user")
	if u == nil {
		resp.Msg = "You are not logged please signin"
		return c.JSON(http.StatusForbidden, resp)
	}

	// get User show
	show, found, err := u.(core.User).GetShowByUUID(c.Param("uuid"))
	if err != nil {
		logger.Log.Error(c.Request().RemoteAddr, " - AjDeleteShow -  ", err)
		return c.JSON(http.StatusInternalServerError, resp)
	}
	if !found {
		resp.Msg = "show not found"
		return c.JSON(http.StatusOK, resp)
	}

	// if locked (sync in progress for ex)
	if show.Locked {
		resp.Msg = "show is locked by other process (sync probably), try again later."
		return c.JSON(http.StatusOK, resp)
	}

	// delete
	logger.Log.Info(c.Request().RemoteAddr, " - ", u.(core.User).Email, " - AjDeleteShow for show ", show.UUID)
	if err = show.Delete(); err != nil {
		logger.Log.Error(c.Request().RemoteAddr, " - ", u.(core.User).Email, " - AjDeleteShow: show.Delete - ", err)
		return c.JSON(http.StatusInternalServerError, resp)
	}
	resp.Ok = true
	return c.JSON(http.StatusOK, resp)
}
