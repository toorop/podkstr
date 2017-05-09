package core

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Episode represents an Show.Episodes
type Episode struct {
	gorm.Model
	ShowID uint
	UUID   string

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

// Enclosure is a Episode.Enclosures
type Enclosure struct {
	EpisodeID uint
	URLimport string
	URL       string
	Length    string
	Type      string
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

// Create creat a nw episode in DB
func (e *Episode) Create() error {
	return DB.Create(e).Error
}
