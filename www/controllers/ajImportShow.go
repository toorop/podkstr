package controllers

import (
	"fmt"
	"net/http"
	"time"

	"strings"

	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"github.com/toorop/podkstr/core"
	"github.com/toorop/podkstr/www/appContext"
	"github.com/toorop/podkstr/www/logger"
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

	// Check if user elready have this feed
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

	//return c.JSON(http.StatusOK, resp)

	feed, err := core.NewFeed(fd.FeedURL)
	if err != nil {
		resp.Msg = err.Error()
		return c.JSON(http.StatusOK, resp)
	}

	// Create show
	image := core.ShowImage{}
	if feed.Channel.Image != nil {
		image = core.ShowImage{
			URL:    feed.Channel.Image.URL,
			Title:  feed.Channel.Image.Title,
			Link:   feed.Channel.Image.Link,
			Width:  feed.Channel.Image.Width,
			Height: feed.Channel.Image.Height,
		}
	}
	show := core.Show{
		UUID:        uuid.NewV4().String(),
		UserID:      u.(core.User).ID,
		Title:       feed.Channel.Title,
		LinkImport:  feed.Channel.Link,
		Link:        feed.Channel.Link,
		Locked:      false,
		Feed:        fd.FeedURL,
		Category:    feed.Channel.Category,
		Description: feed.Channel.Description,
		Subtitle:    feed.Channel.ItunesSubtitle,
		Language:    feed.Channel.Language,
		Copyright:   feed.Channel.Copyright,
		Author:      feed.Channel.ItunesAuthor,
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
		for _, episode := range feed.Channel.Items {
			/*GUIDisPermalink, err := strconv.ParseBool(episode.GUIDisPermalink)
			if err != nil {
				return err
			}*/
			// pubdate
			pubDate, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", episode.PubDate)
			if err != nil {
				return
			}

			// Enclosure
			enclosure := core.Enclosure{
				URLimport: episode.Enclosure.URL,
				Length:    episode.Enclosure.Length,
				Type:      episode.Enclosure.Type,
			}

			// keywords
			keywords := []core.Keyword{}
			ks := strings.Split(episode.ItunesKeywords, ",")
			for _, word := range ks {
				word = strings.ToLower(strings.TrimSpace(word))
				if word != "" {
					// exists ?
					k, found, err := core.GetKeyword(word)
					if err != nil {
						return
					}
					if found {

						keywords = append(keywords, k)
					} else {
						keywords = append(keywords, core.Keyword{
							Word: word,
						})
					}
				}
			}
			ep := core.Episode{
				UUID:            uuid.NewV4().String(),
				Title:           episode.Title,
				LinkImport:      episode.Link,
				Description:     episode.Description,
				Subtitle:        episode.ItunesSubtitle,
				GUID:            episode.GUID,
				GUIDisPermalink: false,
				PubDate:         pubDate,
				Duration:        episode.ItunesDuration,
				Enclosure:       enclosure,
				Keywords:        keywords,

				GoogleplayExplicit: feed.Channel.GoogleplayExplicit,
				ItunesExplicit:     feed.Channel.ItunesExplicit,
			}
			if err := show.AddEpisode(ep); err != nil {
				logger.Log.Error("AjImportShow - show.AddEpisode(ep) ", err)
				return
			}

		}
	}()

	fmt.Println("Au final on a: ", err)

	resp.Ok = true
	return c.JSON(http.StatusOK, resp)
}
