package core

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Show represents a Show. Amazing !!!!!!
type Show struct {
	gorm.Model
	UUID     string `gorm:"type:char(36);unique_index"`
	UserID   uint   `gorm:"index"`
	Locked   bool
	Task     string `gorm:"index"`
	LastSync time.Time

	Title          string `gorm:"type:varchar(1024)"`
	Link           string `gorm:"type:varchar(1024)"`
	LinkImport     string `gorm:"type:varchar(1024)"`
	Feed           string `gorm:"type:varchar(1024)"`
	FeedImport     string `gorm:"type:varchar(1024)"`
	Category       string
	Description    string `gorm:"type:text"`
	Subtitle       string `gorm:"type:text"`
	Language       string
	Copyright      string
	Image          ShowImage
	Author         string
	Owner          string
	OwnerEmail     string
	ItunesExplicit string

	ItunesImage string

	AtomLink string `gorm:"type:varchar(1024)"`

	Episodes []Episode
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
	err = DB.Model(s).Related(&episodes, "Episodes").Error
	return
}

// GetLastEpisode retuns the last episode
func (s *Show) GetLastEpisode() (episode Episode, err error) {
	err = DB.Model(s).Related("Episodes").Order("PubDate").First(&episode).Error
	return
}

// AddEpisode add an episode to the show
func (s *Show) AddEpisode(episode Episode) error {
	s.Episodes = append(s.Episodes, episode)
	return s.Save()
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
