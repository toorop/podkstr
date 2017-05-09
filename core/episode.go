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

	Title           string
	Link            string
	LinkImport      string
	Description     string
	Subtitle        string
	GUID            string
	GUIDisPermalink bool
	PubDate         time.Time
	Duration        uint
	Enclosures      []Enclosure
	Keywords        []Keyword `gorm:"many2many:epidode_keywords"`
	Explicite       string
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
	Episodes []Episode `gorm:"many2many:epidode_keywords"`
}

// Create creat a nw episode in DB
func (e *Episode) Create() error {
	return DB.Create(e).Error
}
