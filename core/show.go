package core

import "github.com/jinzhu/gorm"

// Show represents a Show. Amazing !!!!!!
type Show struct {
	gorm.Model
	UUID   string // hash link
	UserID uint

	Title          string
	Link           string
	LinkImport     string
	Category       string
	Description    string
	Subtitle       string
	Language       string
	Copyright      string
	Image          ShowImage
	Author         string
	ItunesExplicit string
	ItunesOwner    string

	ItunesImage string

	AtomLink string

	Episodes []Episode
}

// ShowImage for Show.Images
type ShowImage struct {
	gorm.Model
	ShowID uint
	URL    string
	Title  string
	Link   string
	Width  string
	Height string
}

// Create nenw show in DB
func (s *Show) Create() error {
	return DB.Create(s).Error
}

// Save saves show in DB
func (s *Show) Save() error {
	return DB.Save(s).Error
}
