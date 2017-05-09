package controllers

import (
	"net/http"

	"strings"

	"fmt"

	"github.com/labstack/echo"
	"github.com/toorop/podkstr/core"
	"github.com/toorop/podkstr/www/appContext"
	"github.com/toorop/podkstr/www/logger"
)

// AjImportShow import
func AjImportShow(ec echo.Context) error {
	type response struct {
		Ok  bool
		Msg string
	}
	var resp = response{}
	var err error

	c := ec.(*appContext.AppContext)
	if c.Get("uEmail").(string) == "" {
		resp.Msg = "You are not logged please signin"
		return c.JSON(http.StatusForbidden, resp)
	}
	type FormData struct {
		FeedURL string `json:"feedURL"`
	}

	fd := new(FormData)
	if err = c.Bind(&fd); err != nil {
		logger.Log.Error(c.Request().RemoteAddr, " - AjImportShow -  ", err)
		return c.JSON(http.StatusInternalServerError, resp)
	}

	fd.FeedURL = strings.TrimSpace(fd.FeedURL)
	if fd.FeedURL == "" {
		resp.Msg = "you must specified a valid feed URL"
		return c.JSON(http.StatusOK, resp)
	}

	feed, err := core.NewFeed(fd.FeedURL)
	if err != nil {
		resp.Msg = err.Error()
		return c.JSON(http.StatusOK, resp)
	}

	for _, episode := range feed.Channel.Item {
		fmt.Println(episode.Title, episode.Link)
	}

	resp.Ok = true
	return c.JSON(http.StatusOK, resp)
}
