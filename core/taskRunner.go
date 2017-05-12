package core

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/h2non/filetype.v1"

	"github.com/spf13/viper"
	"github.com/toorop/podkstr/logger"
)

// TODO handle gracefull shutdown

// TaskRunner . You know what ?! its job is to run tasks...
type TaskRunner struct {
	chanStop chan bool
}

// NewTaskRunner returns a new Taskrunner. Oh Wait ! really ?!!!
func NewTaskRunner() (tr TaskRunner) {
	tr.chanStop = make(chan bool)
	return
}

// Run daemonize TaskRunner 🤘😈🤘
func (tr TaskRunner) Run() {
	for {
		select {
		case _ = <-tr.chanStop:
			logger.Log.Info("TaskRunner - exiting")
			return
		default:
			tr.loop()
		}
	}
}

// Stop stops taskrunner
func (tr TaskRunner) Stop() {
	tr.chanStop <- true
}

// loop is a tasks loop
// for now we will run task sequentially - KISS
func (tr *TaskRunner) loop() {
	logger.Log.Info("TaskRunner - new loop")

	// get show with task to do
	shows, err := GetShowsWithTasks()
	if err != nil {
		logger.Log.Error("TaskRunner - GetShowsWithTasks - ", err)
		return
	}

	for _, show := range shows {
		logger.Log.Info("TaskRunner - begining task ", show.Task, " for show ", show.ID)
		switch show.Task {
		case "sync":
			if err := syncShow(&show); err != nil {
				logger.Log.Error("TaskRunner - syncShow for show ", show.ID, " - ", err)
			}
		case "firstsync":
			if err := firstSyncShow(&show); err != nil {
				logger.Log.Error("TaskRunner -  First syncShow for show ", show.ID, " - ", err)
			}
		default:
			logger.Log.Error("TaskRunner - unknow task ", show.Task, " for show ", show.ID, " - ", err)
		}
		logger.Log.Info("TaskRunner - ending task ", show.Task, " for show ", show.ID)
	}

	// synchro des episodes en backup

	time.Sleep(120 * time.Second)
}

// firstSyncShow firts syn for a show
func firstSyncShow(show *Show) (err error) {
	if err = show.Lock(); err != nil {
		return
	}
	defer show.Unlock()

	// get image show
	image, err := show.GetImage()
	if err != nil {
		return err
	}

	// TODO check if there is an image

	parts := strings.Split(image.URLimport, "/")
	fileName := parts[len(parts)-1]
	// Dl image

	resp0, err := http.Get(image.URLimport)
	if err != nil {
		return err
	}
	defer resp0.Body.Close()
	key := fmt.Sprintf("show/%s/%s", show.UUID, url.QueryEscape(fileName))
	image.StorageKey = key
	image.URL = viper.GetString("openstack.container.url") + "/" + key
	show.ItunesImage = image.URL

	// push to object storage
	filePath := viper.GetString("temppath") + "/image_" + show.UUID
	fd, err := os.Create(filePath)
	if err != nil {
		return err
	}
	_, err = io.Copy(fd, resp0.Body)
	if err != nil {
		fd.Close()
		return err
	}
	fd.Seek(0, 0)
	if err = Store.Put(key, fd); err != nil {
		logger.Log.Error(fmt.Sprintf("TaskRunner - Store.Put(%s) - %s", key, err))
		fd.Close()
		return
	}
	fd.Close()

	// save image
	if err = image.Save(); err != nil {
		logger.Log.Error(fmt.Sprintf("TaskRunner - ShowImage.Save - %s", err))
		return
	}

	// Remove temp file
	os.Remove(filePath)

	// get all episode
	episodes, err := show.GetEpisodes()
	if err != nil {
		return err
	}
	for _, episode := range episodes {

		////////////////////
		// ep Image
		image, found, err := episode.GetImage()
		if err != nil {
			logger.Log.Error(fmt.Sprintf("TaskRunner - episode.GetImage() - %s", err))
			continue
		}
		// TODO Handle error
		logger.Log.Debug("IMAGE", image, image == Image{})

		if found {
			parts := strings.Split(image.URLimport, "/")
			fileName := parts[len(parts)-1]
			// Dl image
			resp1, err := http.Get(image.URLimport)
			if err != nil {
				logger.Log.Error(fmt.Sprintf("TaskRunner - http.Get(%s) - %s", image.URLimport, err))
				continue
			}
			defer resp1.Body.Close()
			key := fmt.Sprintf("show/%s/episode/%s/%s", show.UUID, episode.UUID, url.QueryEscape(fileName))
			image.StorageKey = key
			image.URL = viper.GetString("openstack.container.url") + "/" + key

			// push to object storage
			filePath := viper.GetString("temppath") + "/image_" + episode.UUID
			fd, err := os.Create(filePath)
			if err != nil {
				logger.Log.Error(fmt.Sprintf("TaskRunner - os.Create(%s) - %s", filePath, err))
				continue
			}
			_, err = io.Copy(fd, resp1.Body)
			if err != nil {
				fd.Close()
				logger.Log.Error(fmt.Sprintf("TaskRunner - io.copy - %s", err))
				continue
			}
			fd.Seek(0, 0)
			if err = Store.Put(key, fd); err != nil {
				logger.Log.Error(fmt.Sprintf("TaskRunner - Store.Put(%s) - %s", key, err))
				fd.Close()
				continue
			}
			fd.Close()

			// save image
			if err = image.Save(); err != nil {
				logger.Log.Error(fmt.Sprintf("TaskRunner - Image.Save - %s", err))
				continue
			}

			// Remove temp file
			os.Remove(filePath)
		}

		////////////////////
		// enclosure
		enclosure, err := episode.GetEnclosure()
		if err != nil {
			return err
		}
		resp, err := http.Get(enclosure.URLimport)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		filePath = viper.GetString("temppath") + "/" + episode.UUID
		fd, err = os.Create(filePath)
		if err != nil {
			return err
		}

		written, err := io.Copy(fd, resp.Body)
		fd.Close()
		if err != nil {
			return err
		}

		// set size in DB
		enclosure.Length = written
		if err = enclosure.Update(); err != nil {
			logger.Log.Error("TaskRunner - enclosure.Update - ", err)
		}

		// Get mime type & extension
		buf, _ := ioutil.ReadFile(filePath)

		kind, unkwown := filetype.Match(buf)
		if unkwown == nil {
			logger.Log.Info(fmt.Sprintf("TaskRunner - file %s type: %s - MIME: %s", filePath, kind.Extension, kind.MIME.Value))
			// set duration
			if kind.Extension == "mp3" {
				mp3, err := NewMp3(filePath)
				if err == nil {
					episode.Duration, err = mp3.GetDuration()
					logger.Log.Debug(episode.Duration)
					if err != nil && err != io.EOF {
						logger.Log.Error("TaskRunner - mp3.GetDuration - ", err)
					}
				} else {
					logger.Log.Error("TaskRunner - NewMp3 - ", err)
				}
			}
		}

		// get file name (same as ori)
		parts = strings.Split(enclosure.URLimport, "/")
		fileName = parts[len(parts)-1]

		// push to object storage
		key = fmt.Sprintf("show/%s/episode/%s/%s", show.UUID, episode.UUID, url.QueryEscape(fileName))
		enclosure.StorageKey = key
		logger.Log.Debug("starting transfert for ", key, Store)
		fd, err = os.Open(filePath)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("TaskRunner - os.Open(%s) - %s", filePath, err))
			continue
		}
		if err = Store.Put(key, fd); err != nil {
			logger.Log.Error(fmt.Sprintf("TaskRunner - Store.Put(%s) - %s", key, err))
			continue
		}
		logger.Log.Info("episode transfered to ", key, " - ", err)

		// update enclosure URL
		enclosure.URL = viper.GetString("openstack.container.url") + "/" + key
		if err = enclosure.Update(); err != nil {
			logger.Log.Error("TaskRunner - enclosure.Save - ", err)
		}

		// update episode
		if err = episode.Update(); err != nil {
			logger.Log.Error("TaskRunner - episode.Update - ", err)
		}
		// remove temp file
		os.Remove(filePath)
		break
	}
	// TODO update show status
	show.LastSync = time.Now()
	// TODO do not upadte on err
	show.Task = "sync"
	if err = show.Save(); err != nil {
		logger.Log.Error("TaskRunner - show.Save - ", err)
	}
	return nil
}

// syncShow sync a show (backup)
func syncShow(show *Show) (err error) {

	/*
		feed, err := NewFeed(show.FeedImport)
		if err != nil {
			return err
		}
		lastBuildDate, err = time.Parse(time.RFC1123Z, feed.Channel.LastBuildDate)
		if err != nil {
			// On va chercher pour ce show le dernier episode
			lastBuildDate = time.Now()
		}
		logger.Log.Debug("Feed last build", feed.Channel.LastBuildDate, " - ", lastBuildDate)

		if show.LastSync.Before(lastBuildDate) {
			// new Episodes ?
			logger.Log.Debug("new sync todo")
			lastLocalEpisode, err := show.GetLastEpisode()
			if err != nil {
				return err
			}
			lastSyncDate := lastLocalEpisode.PubDate

			for _, episode := range feed.Channel.Items {
				currentPubDate, err := time.Parse(time.RFC1123Z, episode.PubDate)
				if err != nil {
					return err
				}
				if currentPubDate.After(lastLocalEpisode.PubDate) {
					logger.Log.Debug("Last sync for episode ", episode.PubDate, lastLocalEpisode.PubDate)
				}

			}
		}
	*/

	return nil
}
