package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path"

	"strings"

	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"github.com/toorop/podkstr/core"
	"github.com/toorop/podkstr/logger"
)

// GetShowFeed returns RSS/Atom feed of the show
func GetShowFeed(c echo.Context) error {
	showUUID := c.Param("uuid")

	type tpl struct {
		BaseURL string
		Show    *core.Show
	}

	var err error
	var data = tpl{
		BaseURL: viper.GetString("baseurl"),
	}

	// get feed
	show, found, err := core.GetShowByUUID(showUUID)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("%s - GetShowFeed -> core.GetShowByUUID(%s) - %s ", c.Request().RemoteAddr, showUUID, err))
		return c.NoContent(http.StatusInternalServerError)
	}
	if !found {
		return c.String(http.StatusNotFound, "show not found")
	}

	// if locked or not sync
	if show.Locked || show.Task == "firstsync" {
		return c.String(http.StatusNotFound, "show locked or not synched yet")
	}

	// get episodes
	show.Episodes, err = show.GetEpisodes()
	if err != nil {
		logger.Log.Error(fmt.Sprintf("%s - GetShowFeed -> show.GetEpisodes - %s ", c.Request().RemoteAddr, err))
		return c.NoContent(http.StatusInternalServerError)
	}

	// Enclosures
	for i, ep := range show.Episodes {
		show.Episodes[i].Enclosure, _, err = ep.GetEnclosure()
		if err != nil {
			logger.Log.Error(fmt.Sprintf("%s - GetShowFeed -> ep.GetEnclosure() - %s ", c.Request().RemoteAddr, err))
			return c.NoContent(http.StatusInternalServerError)
		}

		// Enclosure URL must be http and not https (WTF !!!)
		if strings.HasPrefix(show.Episodes[i].Enclosure.URL, "https://") {
			show.Episodes[i].Enclosure.URL = "http://" + show.Episodes[i].Enclosure.URL[8:]
		}
	}

	// Episode Image
	for i, ep := range show.Episodes {
		show.Episodes[i].Image, _, err = ep.GetImage()
		if err != nil {
			logger.Log.Error(fmt.Sprintf("%s - GetShowFeed -> ep.GetImage() - %s ", c.Request().RemoteAddr, err))
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	// Show Image
	show.Image, err = show.GetImage()
	if err != nil {
		logger.Log.Error(fmt.Sprintf("%s - GetShowFeed -> show.GetImage - %s ", c.Request().RemoteAddr, err))
		return c.NoContent(http.StatusInternalServerError)
	}

	data.Show = &show

	// load template
	tplF, err := template.ParseFiles(path.Join(viper.GetString("rootPath"), "etc/tpl/showfeed.rss"))
	if err != nil {
		logger.Log.Error(fmt.Sprintf("%s - GetShowFeed -> template.ParseFiles - %s ", c.Request().RemoteAddr, err))
		return c.NoContent(http.StatusInternalServerError)
	}

	buf := bytes.Buffer{}
	if err = tplF.Execute(&buf, data); err != nil {
		logger.Log.Error(fmt.Sprintf("%s - GetShowFeed -> mailTpl.Execute - %s ", c.Request().RemoteAddr, err))
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.XMLBlob(http.StatusOK, buf.Bytes())
}
