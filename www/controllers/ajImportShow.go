package controllers

import (
	"fmt"
	"net/http"
	"time"

	"strings"

	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"github.com/toorop/podkstr/core"
	"github.com/toorop/podkstr/logger"
	"github.com/toorop/podkstr/www/appContext"
)

// AjImportShow import
func AjImportShow(ec echo.Context) error {
	type response struct {
		Ok   bool
		Msg  string
		Show *core.Show
	}
	var resp = response{}
	var err error

	c := ec.(*appContext.AppContext)
	u := c.Get("user")
	if u == nil {
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

	// Check if user already have this feed
	_, found, err := u.(core.User).GetShowByFeed(fd.FeedURL)
	if err != nil {
		resp.Msg = err.Error()
		return c.JSON(http.StatusOK, resp)
	}
	if found {
		resp.Ok = false
		resp.Msg = "You already have this show on your show list."
		return c.JSON(http.StatusOK, resp)
	}

	feed, err := core.NewFeed(fd.FeedURL)
	if err != nil {
		resp.Msg = err.Error()
		return c.JSON(http.StatusOK, resp)
	}

	// Create show

	// Last build date
	lastBuildDate, err := time.Parse(time.RFC1123Z, feed.Channel.LastBuildDate)
	if err != nil {
		lastBuildDate = time.Now()
	}

	// image
	image := core.ShowImage{}
	if feed.Channel.Image != (core.FeedImage{}) {
		image = core.ShowImage{
			URL:       feed.Channel.Image.URL,
			URLimport: feed.Channel.Image.URL,
			Title:     feed.Channel.Image.Title,
			Link:      feed.Channel.Image.Link,
			Width:     feed.Channel.Image.Width,
			Height:    feed.Channel.Image.Height,
		}
	}

	show := core.Show{
		UUID:          uuid.NewV4().String(),
		Locked:        false,
		Task:          "firstsync",
		UserID:        u.(core.User).ID,
		Title:         feed.Channel.Title,
		LastBuildDate: lastBuildDate,
		LastSync:      time.Now(),

		LinkImport:  feed.Channel.Link,
		Link:        feed.Channel.Link,
		AtomLink:    feed.Channel.AtomLink.Href,
		Feed:        fd.FeedURL,
		FeedImport:  fd.FeedURL,
		Category:    feed.Channel.Category,
		Description: feed.Channel.Description,
		Subtitle:    feed.Channel.ItunesSubtitle,
		Language:    feed.Channel.Language,
		Copyright:   feed.Channel.Copyright,
		Author:      feed.Channel.ItunesAuthor,
		Owner:       feed.Channel.ItunesOwner.Name,
		OwnerEmail:  feed.Channel.ItunesOwner.Email,
		Image:       image,

		ItunesExplicit: feed.Channel.ItunesExplicit,
		ItunesImage:    feed.Channel.ItunesImage.Href,

		Episodes: []core.Episode{},
	}

	if err = show.Create(); err != nil {
		return err
	}

	resp.Show = &show

	go func() {
		for _, feedEpisode := range feed.Channel.Items {
			_, err := show.AddEpisodeFromFeed(feedEpisode)
			if err != nil {
				logger.Log.Error("AjImportShow - show.AddEpisodeFromFeed(ep) ", err)
				return
			}
		}
	}()

	fmt.Println("Au final on a: ", err)

	resp.Ok = true
	return c.JSON(http.StatusOK, resp)
}
