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

	Title              string
	Link               string
	LinkImport         string
	Description        string
	Subtitle           string
	GUID               string
	GUIDisPermalink    bool
	PubDate            time.Time
	Duration           string
	Enclosure          Enclosure
	Keywords           []Keyword `gorm:"many2many:episode_keywords"`
	ItunesExplicit     string
	GoogleplayExplicit string
}

// Create creat a nw episode in DB
func (e *Episode) Create() error {
	return DB.Create(e).Error
}

// Delete delete an episode
func (e Episode) Delete() (err error) {

	// Enclosure
	// Get episode enclosure
	var enclosure Enclosure
	if err = DB.Model(&e).Related(&enclosure).Error; err != nil {
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

// GetKeywords returns episode keywords
func (e *Episode) GetKeywords() (keywords []Keyword, err error) {
	err = DB.Model(e).Related(&keywords, "Keywords").Error
	return
}

// Enclosure is a Episode.Enclosures
type Enclosure struct {
	gorm.Model
	EpisodeID uint `gorm:"index"`
	URLimport string
	URL       string
	Length    string
	Type      string
}

// Delete delete enclosure e
func (e *Enclosure) Delete() error {
	// TODO delete file
	return DB.Unscoped().Delete(e).Error
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
