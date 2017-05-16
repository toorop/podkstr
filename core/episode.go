package core

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Episode represents an Show.Episodes
type Episode struct {
	gorm.Model
	ShowID uint   `gorm:"index"`
	UUID   string `gorm:"type:char(36);unique_index"`

	Title              string `gorm:"type:varchar(1024)"`
	Link               string `gorm:"type:varchar(1024)"`
	LinkImport         string `gorm:"type:varchar(1024)"`
	Description        string `gorm:"type:text"`
	Subtitle           string `gorm:"type:text"`
	GUID               string
	GUIDisPermalink    bool
	PubDate            time.Time
	Duration           time.Duration
	Image              Image
	Enclosure          Enclosure
	Keywords           []Keyword `gorm:"many2many:episode_keywords"`
	ItunesExplicit     string
	GoogleplayExplicit string
}

// Create creat a nw episode in DB
func (e *Episode) Create() error {
	return DB.Create(e).Error
}

// Update update episode
func (e *Episode) Update() error {
	return DB.Save(e).Error
}

// Delete delete an episode
func (e Episode) Delete() (err error) {

	// Image
	var image Image
	image, found, err := e.GetImage()
	if err != nil {
		return err
	}
	if found {
		if err = image.Delete(); err != nil {
			return err
		}
	}
	// Enclosure
	var enclosure Enclosure
	enclosure, err = e.GetEnclosure()
	if err != nil {
		return err
	}
	if err = enclosure.Delete(); err != nil {
		return err
	}

	// delete episode keywords
	/*var keywords []Keyword
	if err = DB.Unscoped().Model(e).Related(&keywords, "Keywords").Delete(&keywords).Error; err != nil {
		return err
	}*/

	// Pour le moment on ne supprime que les associations
	if err = DB.Model(e).Association("Keywords").Clear().Error; err != nil {
		return err
	}

	// delete episode from DB
	return DB.Unscoped().Delete(e).Error
}

// GetImage return episode image
func (e *Episode) GetImage() (image Image, found bool, err error) {
	err = DB.Model(e).Related(&image).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		return
	}
	found = true
	return
}

// GetKeywords returns episode keywords
func (e *Episode) GetKeywords() (keywords []Keyword, err error) {
	err = DB.Model(e).Related(&keywords, "Keywords").Error
	return
}

// GetEnclosure return episode enclosure
func (e *Episode) GetEnclosure() (enclosure Enclosure, err error) {
	err = DB.Model(e).Related(&enclosure).Error
	return
}

// Enclosure is a Episode.Enclosures
type Enclosure struct {
	gorm.Model
	EpisodeID  uint `gorm:"index"`
	URLimport  string
	URL        string
	StorageKey string
	Length     int64
	Type       string
}

// Delete delete enclosure e
func (e *Enclosure) Delete() error {
	// delete from storage
	if e.StorageKey != "" {
		if err := Store.Del(e.StorageKey); err != nil {
			return err
		}
	}
	return DB.Unscoped().Delete(e).Error
}

// Update update enclosure
func (e *Enclosure) Update() error {
	return DB.Save(e).Error
}

// Image represents an Episode.Image
type Image struct {
	gorm.Model
	EpisodeID  uint `gorm:"index"`
	URL        string
	URLimport  string
	Title      string
	Link       string
	LinkImport string
	StorageKey string
}

// Delete delete an image
func (i *Image) Delete() error {
	// delete from storage
	if i.StorageKey != "" {
		if err := Store.Del(i.StorageKey); err != nil {
			return err
		}
	}
	return DB.Unscoped().Delete(i).Error
}

// Save update Image
func (i *Image) Save() error {
	return DB.Save(i).Error
}

// Keyword is a Episode.Keywords
type Keyword struct {
	gorm.Model
	Word     string
	Episodes []Episode `gorm:"many2many:episode_keywords"`
}

// GetKeyword return Keyword
func GetKeyword(word string) (k Keyword, found bool, err error) {
	err = DB.Where("word = ?", word).First(&k).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		return
	}
	found = true
	return
}
