package core

// TODO refactor this sh*t

import (
	"fmt"
	"time"

	"github.com/satori/go.uuid"
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
		if show.Task == "" {
			continue
		}
		taskID := uuid.NewV4().String()
		logger.Log.Info(fmt.Sprintf("TaskRunner - begining task %s - %s for show %d", taskID, show.Task, show.ID))
		switch show.Task {
		case "sync":
			if err := syncShow(&show, taskID); err != nil {
				logger.Log.Error(fmt.Sprintf("TaskRunner - task %s - %s", taskID, err))
			}
		case "firstsync":
			if err := firstSyncShow(&show); err != nil {
				logger.Log.Error(fmt.Sprintf("TaskRunner - task %s - %s", taskID, err))
			}
		default:
			logger.Log.Error(fmt.Sprintf("TaskRunner - task %s - uknow task %s", taskID, show.Task))
		}
		logger.Log.Info(fmt.Sprintf("TaskRunner - end task %s - %s for show %d", taskID, show.Task, show.ID))

	}
	logger.Log.Info("TaskRunner - end tasks loop")
	time.Sleep(120 * time.Second)
}

// firstSyncShow firts syn for a show
func firstSyncShow(show *Show) (err error) {
	if err = show.Lock(); err != nil {
		return
	}
	defer show.Unlock()

	// get show image
	image, found, err := show.GetImage()
	if err != nil {
		return err
	}

	// if there is an image
	if found {
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
	show.Task = "sync"
	if err = show.Save(); err != nil {
		logger.Log.Error("TaskRunner - show.Save - ", err)
	}
	return nil
}

// syncShow sync a show (backup)
func syncShow(show *Show, taskID string) (err error) {
	if err = show.Lock(); err != nil {
		return
	}
	defer show.Unlock()
	feed, err := NewFeed(show.FeedImport)
	if err != nil {
		return err
	}

	// Pour chaque épisode du feed on vérifie qu'il existe en DB
	for _, feedEpisode := range feed.Channel.Items {
		if feedEpisode.GUID == "" {
			return fmt.Errorf("syncShow - feedEpisode has no GUID")
		}
		episode, found, err := show.GetEpisodeByGUID(feedEpisode.GUID)
		if err != nil {
			return fmt.Errorf("syncShow GetEpisodeFromFeedEpisode - %s", err)
		}
		//logger.Log.Debugf("EPISODE BU GUID %v - Found %v", episode, found)
		//continue

		// if episode not found sync
		if !found {
			episode, err = show.AddEpisodeFromFeed(feedEpisode)

			if err != nil {
				return fmt.Errorf("syncShow episode.Sync() for show %d, episode %d(%s) - %s", show.ID, episode.ID, episode.UUID, err)
			}
			if err = episode.Sync(); err != nil {
				return fmt.Errorf("syncShow episode.Sync() for show %d, episode %d(%s) - %s", show.ID, episode.ID, episode.UUID, err)
			}
			logger.Log.Info(fmt.Sprintf("TaskRunner - task %s for show %d, episode %d (%s) - resync succesfully done", taskID, episode.ShowID, episode.ID, episode.UUID))

			continue

		} else {
			// if episode found check data && synch

			// image
			image, found, err := episode.GetImage()
			if err != nil {
				return fmt.Errorf("syncShow episode.GetImage() - %s", err)
			}

			// if feed episode has image but not podkstr copy
			if (feedEpisode.Image != (ItemImage{})) && !found {
				logger.Log.Info(fmt.Sprintf("TaskRunner - task %s for show %d, episode %d (%s) - image is not normally sync, processing resync", taskID, episode.ShowID, episode.ID, episode.UUID))
				if err = episode.Sync(); err != nil {
					return fmt.Errorf("syncShow episode.Sync() (image) for show %d, episode %d(%s) - %s", episode.ShowID, episode.ID, episode.UUID, err)
				}
				logger.Log.Info(fmt.Sprintf("TaskRunner - task %s for show %d, episode %d (%s) - resync (image) succesfully done", taskID, episode.ShowID, episode.ID, episode.UUID))
				continue
			}

			// check sync
			if found {
				if image.URL == "" || image.StorageKey == "" {
					logger.Log.Info(fmt.Sprintf("TaskRunner - task %s for show %d, episode %d (%s) - image is not normally sync, processing resync", taskID, episode.ShowID, episode.ID, episode.UUID))
					if err = episode.Sync(); err != nil {
						return fmt.Errorf("syncShow episode.Sync() (image) for show %d, episode %d(%s) - %s", episode.ShowID, episode.ID, episode.UUID, err)
					}
					logger.Log.Info(fmt.Sprintf("TaskRunner - task %s for show %d, episode %d (%s) - resync (image) succesfully done", taskID, episode.ShowID, episode.ID, episode.UUID))
					continue
				}
			}

			// enclosure
			enclosure, found, err := episode.GetEnclosure()

			//logger.Log.Debugf("ENCLOSURES: %v", enclosure)

			// If feed has enclosure but local copy don't
			if (feedEpisode.Enclosure != (ItemEnclosure{})) && !found {
				logger.Log.Info(fmt.Sprintf("TaskRunner - task %s for show %d, episode %d (%s) - enclosure not found locally, processing resync", taskID, episode.ShowID, episode.ID, episode.UUID))
				if err = episode.Sync(); err != nil {
					return fmt.Errorf("syncShow episode.Sync() (enclosure) for show %d, episode %d(%s) - %s", episode.ShowID, episode.ID, episode.UUID, err)
				}
				logger.Log.Info(fmt.Sprintf("TaskRunner - task %s for show %d, episode %d (%s) - resync (enclosure) succesfully done", taskID, episode.ShowID, episode.ID, episode.UUID))
				continue
			}

			// check if sync is successfully done
			if found {
				//logger.Log.Debugf("ENCLOSURE URL |%s| - ENCLOSURE KEY |%s|", enclosure.URL, enclosure.StorageKey)
				if enclosure.URL == "" || enclosure.StorageKey == "" {
					logger.Log.Info(fmt.Sprintf("TaskRunner - task %s for show %d, episode %d (%s) - enclosure is not normally sync, processing resync", taskID, episode.ShowID, episode.ID, episode.UUID))
					if err = episode.Sync(); err != nil {
						return fmt.Errorf("syncShow episode.Sync() (enclosure) for show %d, episode %d(%s) - %s", episode.ShowID, episode.ID, episode.UUID, err)
					}
					logger.Log.Info(fmt.Sprintf("TaskRunner - task %s for show %d, episode %d (%s) - resync (enclosure) succesfully done", taskID, episode.ShowID, episode.ID, episode.UUID))
				}
			}
		}
	}

	/*
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
		}*/
	show.LastSync = time.Now()
	// TODO do not update on err
	show.Task = "sync"
	if err = show.Save(); err != nil {
		logger.Log.Error("TaskRunner - show.Save - ", err)
	}
	return nil
}
