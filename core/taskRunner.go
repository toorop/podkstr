package core

// TODO refactor this sh*t

import (
	"fmt"
	"time"

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

	// get show image
	image, err := show.GetImage()
	if err != nil {
		return err
	}

	// TODO check if there is an image
	if image != (ShowImage{}) {
		image.StorageKey, image.URL, err = StoreCopyImageFromURL(fmt.Sprintf("show/%s", show.UUID), image.URLimport)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("TaskRunner - StoreCopyImageFromURL - %s", err))
			return err
		}

		// save image
		if err = image.Save(); err != nil {
			logger.Log.Error(fmt.Sprintf("TaskRunner - ShowImage.Save - %s", err))
			return err
		}
	}

	// itunes image
	if show.ItunesImage != "" {
		_, iURL, err := StoreCopyImageFromURL(fmt.Sprintf("show/%s", show.UUID), show.ItunesImage)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("TaskRunner - StoreCopyImageFromURL - %s", err))
			return err
		}
		show.ItunesImage = iURL
		logger.Log.Debug("SHOW URL ", iURL)
	}

	// get all episode
	episodes, err := show.GetEpisodes()
	if err != nil {
		return err
	}
	for _, episode := range episodes {
		if err = episode.Sync(); err != nil {
			logger.Log.Error(fmt.Sprintf("TaskRunner - %s", err))
			return
		}

	}
	show.LastSync = time.Now()
	// TODO do not update on err
	show.Task = "sync"

	if err = show.Save(); err != nil {
		logger.Log.Error("TaskRunner - show.Save - ", err)
	}
	return nil
}

// syncShow sync a show (backup)
func syncShow(show *Show) (err error) {
	if err = show.Lock(); err != nil {
		return
	}
	defer show.Unlock()
	feed, err := NewFeed(show.FeedImport)
	if err != nil {
		return err
	}
	lastBuildDate, err := time.Parse(time.RFC1123Z, feed.Channel.LastBuildDate)
	if err != nil {
		// On va chercher pour ce show le dernier episode
		lastBuildDate = time.Now()
	}
	logger.Log.Debug("Feed last build from feed ", lastBuildDate, " in DB ", show.LastSync)

	if show.LastSync.Before(lastBuildDate) {
		// new Episodes ?
		logger.Log.Debug("new sync todo")
		lastLocalEpisode, err := show.GetLastEpisode()
		if err != nil {
			return err
		}
		for _, feedEpisode := range feed.Channel.Items {
			fromFeedPubDate, err := time.Parse(time.RFC1123Z, feedEpisode.PubDate)
			if err != nil {
				return err
			}
			if fromFeedPubDate.After(lastLocalEpisode.PubDate) {
				logger.Log.Debug("New episode to sync ", show.Title, feedEpisode.Title)

				// add Episode
				var episodeUUID string
				if episodeUUID, err = show.AddEpisodeFromFeed(feedEpisode); err != nil {
					logger.Log.Error("TaskRunner - show.AddEpisodeFromFeed - ", err)
					return err
				}
				// sync episode
				logger.Log.Debug(fmt.Sprintf("TaskRunner - episode added to DB with UUID - %s", episodeUUID))
				episode, found, err := GetEpisodeByUUID(episodeUUID)
				if err != nil {
					logger.Log.Error(fmt.Sprintf("TaskRunner - GetEpisodeByUUID(%s) - %s ", episodeUUID, err))
					return err
				}
				if !found {
					logger.Log.Error(fmt.Sprintf("TaskRunner - GetEpisodeByUUID(%s) - not found ", episodeUUID))
					return err
				}

				if err = episode.Sync(); err != nil {
					logger.Log.Error(fmt.Sprintf("TaskRunner - episode.Sync() for %s - %s ", episodeUUID, err))
					return err
				}
			}
		}
	}
	show.LastSync = time.Now()
	// TODO do not update on err
	show.Task = "sync"
	if err = show.Save(); err != nil {
		logger.Log.Error("TaskRunner - show.Save - ", err)
	}
	return nil
}
