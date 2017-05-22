package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
)

// Show represents a Show. Amazing !!!!!!
type Show struct {
	gorm.Model
	UUID     string `gorm:"type:char(36);unique_index"`
	UserID   uint   `gorm:"index"`
	Locked   bool
	Task     string `gorm:"index"`
	LastSync time.Time

	Title         string `gorm:"type:varchar(1024)"`
	Link          string `gorm:"type:varchar(1024)"`
	LinkImport    string `gorm:"type:varchar(1024)"`
	Feed          string `gorm:"type:varchar(1024)"`
	FeedImport    string `gorm:"type:varchar(1024)"`
	Category      string
	LastBuildDate time.Time
	Description   string `gorm:"type:text"`
	Subtitle      string `gorm:"type:text"`
	Language      string
	Copyright     string
	Image         ShowImage
	Author        string
	Owner         string
	OwnerEmail    string

	ItunesCategory string
	ItunesExplicit string
	ItunesImage    string

	GoogleplayExplicit string

	AtomLink string `gorm:"type:varchar(1024)"`

	Episodes []Episode
}

// GetShowByID returns show by ID
func GetShowByID(ID uint) (show Show, found bool, err error) {
	err = DB.First(&show, ID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		return
	}
	found = true
	return
}

// GetShowByUUID returns a show by its UUID
func GetShowByUUID(UUID string) (show Show, found bool, err error) {
	err = DB.Where("uuid = ?", UUID).First(&show).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		return
	}
	found = true
	return
}

// GetShowsWithTasks returns show which have tasks to do
func GetShowsWithTasks() (shows []Show, err error) {
	err = DB.Where("task != ''").Find(&shows).Error
	return
}

// Lock set locked flag to true
func (s *Show) Lock() error {
	s.Locked = true
	return s.Save()
}

// Unlock unset locked flag to true
func (s *Show) Unlock() error {
	s.Locked = false
	return s.Save()
}

// Create nenw show in DB
func (s *Show) Create() error {
	return DB.Create(s).Error
}

// Save saves show in DB
func (s *Show) Save() error {
	return DB.Save(s).Error
}

// Update updates show in DB
func (s *Show) Update() error {
	return DB.Save(s).Error
}

// Delete delete show and episodes
func (s *Show) Delete() (err error) {
	if err = s.Lock(); err != nil {
		return err
	}

	defer s.Unlock()

	// delete showImage
	image, err := s.GetImage()
	if err != nil {
		return
	}

	if err = image.Delete(); err != nil {
		return err
	}
	/*if err = DB.Unscoped().Model(s).Related(&ShowImage{}).Delete(&ShowImage{}).Error; err != nil && err != gorm.ErrRecordNotFound {
		s.Unlock()
		return err
	}*/

	// TODO delete ItunesImage
	if s.ItunesImage != "" && strings.HasPrefix(s.ItunesImage, viper.GetString("openstack.container.url")) {
		key := s.ItunesImage[len(viper.GetString("openstack.container.url"))+1:]
		if err = Store.Del(key); err != nil {
			return err
		}

	}

	// Delete épisode
	// get show episodes
	episodes, err := s.GetEpisodes()
	if err != nil {
		return err
	}

	for _, episode := range episodes {
		if err = episode.Delete(); err != nil {
			return err
		}
	}
	// Delete show
	return DB.Unscoped().Delete(s).Error
}

// GetEpisodes return show episodes
func (s *Show) GetEpisodes() (episodes []Episode, err error) {
	//err = DB.Model(s).Related(&Episode{}).Order("pub_date desc").Find(&episodes).Error
	//err = DB.Model(s).Related(&episodes, "Episodes").Order("pub_date desc").Error
	err = DB.Model(&Episode{}).Where("show_id = ?", s.ID).Order("pub_date desc").Find(&episodes).Error

	return
}

// GetLastEpisode retuns the last episode
func (s *Show) GetLastEpisode() (episode Episode, err error) {
	err = DB.Model(s).Related(&Episode{}).Order("pub_date desc").First(&episode).Error
	return
}

// GetEpisodeByGUID returns Show episode by it GUID
func (s *Show) GetEpisodeByGUID(GUID string) (episode Episode, found bool, err error) {
	err = DB.Model(&Episode{}).Where("guid = ? AND show_id = ?", GUID, s.ID).First(&episode).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		return
	}
	found = true
	return
}

// AddEpisode add an episode to the show
func (s *Show) AddEpisode(episode *Episode) error {
	s.Episodes = append(s.Episodes, *episode)
	err := s.Save()
	return err
}

// AddEpisodeFromFeed add an épisode from feed
func (s *Show) AddEpisodeFromFeed(feedEpisode Item) (episode Episode, err error) {
	pubDate, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", feedEpisode.PubDate)
	if err != nil {
		return
	}

	// Image
	image := Image{}
	if feedEpisode.Image != (ItemImage{}) {
		image = Image{
			URL:        feedEpisode.Image.URL,
			URLimport:  feedEpisode.Image.URL,
			Link:       feedEpisode.Image.Link,
			LinkImport: feedEpisode.Image.Link,
			Title:      feedEpisode.Image.Title,
		}
	}

	// Enclosure
	lenght, _ := strconv.ParseInt(feedEpisode.Enclosure.Length, 10, 64)

	enclosure := Enclosure{
		URLimport: feedEpisode.Enclosure.URL,
		Length:    lenght,
		Type:      feedEpisode.Enclosure.Type,
	}

	// keywords
	keywords := []Keyword{}
	ks := strings.Split(feedEpisode.ItunesKeywords, ",")
	for _, word := range ks {
		word = strings.ToLower(strings.TrimSpace(word))
		if word != "" {
			// exists ?
			k, found, err := GetKeyword(word)
			if err != nil {
				return episode, err
			}
			if found {

				keywords = append(keywords, k)
			} else {
				keywords = append(keywords, Keyword{
					Word: word,
				})
			}
		}
	}

	// duration
	// soit du type H:M:S soit en secondes sinon 0 et on le calculera au first sync
	var duration time.Duration
	durationParts := strings.Split(feedEpisode.ItunesDuration, ":")
	if len(durationParts) > 1 {
		if len(durationParts) == 2 {
			duration, _ = time.ParseDuration(fmt.Sprintf("%sm%ss", durationParts[0], durationParts[1]))
		} else if len(durationParts) == 3 {
			duration, _ = time.ParseDuration(fmt.Sprintf("%sh%sm%ss", durationParts[0], durationParts[1], durationParts[2]))
		}
	}
	// second ?
	if duration == 0 {
		duration, _ = time.ParseDuration(feedEpisode.ItunesDuration + "s")
	}

	episode = Episode{
		UUID:            uuid.NewV4().String(),
		Title:           feedEpisode.Title,
		LinkImport:      feedEpisode.Link,
		Description:     feedEpisode.Description,
		Subtitle:        feedEpisode.ItunesSubtitle,
		GUID:            feedEpisode.GUID,
		GUIDisPermalink: false,
		Author:          feedEpisode.ItunesAuthor,
		PubDate:         pubDate,
		Duration:        duration,
		Image:           image,
		Enclosure:       enclosure,
		Keywords:        keywords,

		GoogleplayExplicit: feedEpisode.GoogleplayExplicit,
		ItunesExplicit:     feedEpisode.ItunesExplicit,
	}
	if err = s.AddEpisode(&episode); err != nil {
		return episode, err
	}
	episode, found, err := GetEpisodeByUUID(episode.UUID)
	if err == nil && !found {
		err = fmt.Errorf("episode %d-%s not found", s.ID, episode.UUID)
	}
	return episode, err
}

// GetImage return show image
func (s *Show) GetImage() (image ShowImage, err error) {
	err = DB.Model(s).Related(&image).Error
	return
}

// ShowImage for Show.Images
type ShowImage struct {
	gorm.Model
	ShowID     uint `gorm:"index"`
	URL        string
	URLimport  string
	StorageKey string
	Title      string
	Link       string
	Width      string
	Height     string
}

// Delete delete an image
func (i *ShowImage) Delete() error {
	// delete from storage
	if i.StorageKey != "" {
		if err := Store.Del(i.StorageKey); err != nil {
			return err
		}
	}
	return DB.Unscoped().Delete(i).Error
}

// Save update Image
func (i *ShowImage) Save() error {
	return DB.Save(i).Error
}
