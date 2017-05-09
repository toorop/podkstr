package core

import "github.com/jinzhu/gorm"

// Show represents a Show. Amazing !!!!!!
type Show struct {
	gorm.Model
	UUID  string // hash link
	Owner User

	Title       string
	Link        string
	LinkImport  string
	Category    string
	Description string
	Subtitle    string
	Language    string
	Copyright   string
	Images      []ShowImage
	Author      string
	explicit    bool
	ItunesImage string

	AtomLink string

	Episode []Episode
}

// ShowImage for Show.Images
type ShowImage struct {
	gorm.Model
	ShowID uint
	URL    string
}

// Create nenw show in DB
func (s *Show) Create() error {
	return DB.Create(s).Error
}
