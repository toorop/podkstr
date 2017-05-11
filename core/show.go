package core

import "github.com/jinzhu/gorm"

// Show represents a Show. Amazing !!!!!!
type Show struct {
	gorm.Model
	UUID   string `gorm:"type:char(36);unique_index"`
	UserID uint   `gorm:"index"`
	Locked bool

	Title          string `gorm:"type:varchar(1024)"`
	Link           string `gorm:"type:varchar(1024)"`
	LinkImport     string `gorm:"type:varchar(1024)"`
	Feed           string `gorm:"type:varchar(1024)"`
	Category       string
	Description    string `gorm:"type:text"`
	Subtitle       string `gorm:"type:text"`
	Language       string
	Copyright      string
	Image          ShowImage
	Author         string
	ItunesExplicit string
	ItunesOwner    string

	ItunesImage string

	AtomLink string `gorm:"type:varchar(1024)"`

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
	return DB.Update(s).Error
}

// Delete delete show and episodes
func (s *Show) Delete() (err error) {
	if err = s.Lock(); err != nil {
		return err
	}

	// delete showImage
	if err = DB.Unscoped().Model(s).Related(&ShowImage{}).Delete(&ShowImage{}).Error; err != nil && err != gorm.ErrRecordNotFound {
		s.Unlock()
		return err
	}

	// Delete épisode
	// get show episodes
	episodes, err := s.GetEpisodes()
	if err != nil {
		s.Unlock()
		return err
	}

	for _, episode := range episodes {
		if err = episode.Delete(); err != nil {
			s.Unlock()
			return err
		}

	}
	// Delete show
	if err = DB.Unscoped().Delete(s).Error; err != nil {
		s.Unlock()
	}
	return err
}

// GetEpisodes return show episodes
func (s *Show) GetEpisodes() (episodes []Episode, err error) {
	err = DB.Model(s).Related(&episodes, "Episodes").Error
	return
}

// AddEpisode add an episode to the show
func (s *Show) AddEpisode(episode Episode) error {
	s.Episodes = append(s.Episodes, episode)
	return s.Save()
}
