package core

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/h2non/filetype.v1"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
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
	GUID               string `gorm:"type:varchar(1024);index"` // Original UUID for import
	GUIDisPermalink    bool
	Author             string `gorm:"type:varchar(1024)"`
	PubDate            time.Time
	Duration           time.Duration
	Image              Image
	Enclosure          Enclosure
	Keywords           []Keyword `gorm:"many2many:episode_keywords"`
	ItunesExplicit     string
	GoogleplayExplicit string
}

// GetEpisodeByUUID returns (or not) an episode by UUID
func GetEpisodeByUUID(UUID string) (episode Episode, found bool, err error) {
	err = DB.Where("uuid = ?", UUID).First(&episode).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		return
	}
	found = true
	return
}

// GetEpisodeByGUID returns an episode by its GUID
func GetEpisodeByGUID(GUID string) (episode Episode, found bool, err error) {
	err = DB.Where("guid = ?", GUID).First(&episode).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		return
	}
	found = true
	return
}

// Create creat a new episode in DB
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
	enclosure, found, err = e.GetEnclosure()
	if err != nil {
		return err
	}
	if found {
		if err = enclosure.Delete(); err != nil {
			return err
		}
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

// Sync synchronize an episode
func (e *Episode) Sync() error {
	show, found, err := GetShowByID(e.ShowID)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("show with ID %d not found", e.ShowID)
	}

	////////////////////
	// image
	image, found, err := e.GetImage()
	if err != nil {
		return fmt.Errorf("unable to get image for episode %d - %s", e.ID, err)
	}
	if found {
		image.StorageKey, image.URL, err = StoreCopyImageFromURL(fmt.Sprintf("show/%s", show.UUID), image.URLimport)
		if err != nil {
			return fmt.Errorf("unable to StoreCopyImageFromURL for episode %d - %s", e.ID, err)
		}
		// save image
		if err = image.Save(); err != nil {
			return fmt.Errorf("unable to save image for episode %d - %s", e.ID, err)
		}
	}

	////////////////////
	// enclosure
	enclosure, found, err := e.GetEnclosure()
	if err != nil {
		return fmt.Errorf("unable to getEnclosure() for episode %d - %s", e.ID, err)
	}
	if !found {
		return fmt.Errorf("enclosure not found for episode %d - %s", e.ID, err)
	}

	resp, err := http.Get(enclosure.URLimport)
	if err != nil {
		return fmt.Errorf("unable to http.Get(%s) for episode %d - %s", enclosure.URLimport, e.ID, err)
	}
	defer resp.Body.Close()
	filePath := viper.GetString("temppath") + "/" + e.UUID
	fd, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("unable to os.Create(%s) for episode %d - %s", filePath, e.ID, err)
	}

	written, err := io.Copy(fd, resp.Body)
	fd.Close()
	if err != nil {
		return fmt.Errorf("unable to io.copy() for episode %d - %s", e.ID, err)
	}

	// set size in DB
	enclosure.Length = written
	if err = enclosure.Update(); err != nil {
		return fmt.Errorf("unable to enclosure.Update() for episode %d - %s", e.ID, err)
	}

	// Get mime type & extension
	buf, _ := ioutil.ReadFile(filePath)

	kind, unkwown := filetype.Match(buf)
	if unkwown == nil {
		// set duration
		if kind.Extension == "mp3" {
			mp3, err := NewMp3(filePath)
			if err == nil {
				e.Duration, _ = mp3.GetDuration()
			}
		}
	}

	// get file name (same as ori)
	parts := strings.Split(enclosure.URLimport, "/")
	fileName := parts[len(parts)-1]

	// push to object storage if not found
	// get file hash
	enclosure.Hash, err = GetSHA256File(filePath)
	if err != nil {
		return fmt.Errorf("unable to get SHA256 of %s - %s", filePath, err)
	}
	enc, found, err := GetEnclosureByHash(enclosure.Hash)
	if err != nil {
		return fmt.Errorf("unable to GetEnclosureByHash(%s) - %s", filePath, err)
	}
	if found {
		enclosure.URL = enc.URL
		enclosure.StorageKey = enc.StorageKey
	} else {
		key := fmt.Sprintf("enclosures/%s/%s", enclosure.Hash, url.QueryEscape(fileName))
		enclosure.StorageKey = key
		fd, err = os.Open(filePath)
		if err != nil {
			return fmt.Errorf("unable to os.Open(%s) for episode %d - %s", filePath, e.ID, err)
		}
		defer os.Remove(filePath)
		if err = Store.Put(key, fd); err != nil {
			return fmt.Errorf("unable to store.Put(%s) for episode %d - %s", key, e.ID, err)
		}

		// update enclosure URL
		enclosure.URL = viper.GetString("openstack.container.url") + "/" + key

	}
	// update enclosure
	if err = enclosure.Update(); err != nil {
		return fmt.Errorf("unable to enclosure.Update() for episode %d - %s", e.ID, err)
	}
	// update episode
	if err = e.Update(); err != nil {
		return fmt.Errorf("unable to e.Update() for episode %d - %s", e.ID, err)
	}

	return nil
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
func (e *Episode) GetEnclosure() (enclosure Enclosure, found bool, err error) {
	err = DB.Model(e).Related(&enclosure).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		return
	}
	found = true
	return
}

// FormattedPubDate returns pubdate formatted as String RFC1123Z
func (e *Episode) FormattedPubDate() string {
	return e.PubDate.Format(time.RFC1123Z)
}

// FormatedKeywords returns formated keywords
func (e *Episode) FormatedKeywords() (fKeywords string) {
	kw, err := e.GetKeywords()
	if err != nil {
		return ""
	}
	if len(kw) == 0 {
		return
	}
	for _, w := range kw {
		fKeywords += "," + w.Word
	}
	return fKeywords[1:]
}

// FormattedDuration returns formated duration (for RSS)
func (e *Episode) FormattedDuration() string {
	return fmt.Sprintf("%d", int(e.Duration.Seconds()))
}

// FormattedItunesExplicit returns RSS formated itunes explicit
func (e *Episode) FormattedItunesExplicit() string {
	if e.ItunesExplicit == "" {
		e.ItunesExplicit = "no"
	}
	return e.ItunesExplicit
}

/////////////////////////////////
// Enclosure

// Enclosure is a Episode.Enclosure
type Enclosure struct {
	gorm.Model
	EpisodeID  uint   `gorm:"index"`
	Hash       string `gorm:"type:char(64);index"`
	URLimport  string `gorm:"type:varchar(1024)"`
	URL        string `gorm:"type:varchar(1024)"`
	StorageKey string
	Length     int64
	Type       string
}

// GetEnclosureByHash return the first enclosure with this hash
func GetEnclosureByHash(hash string) (e Enclosure, found bool, err error) {
	err = DB.Where("hash = ?", hash).First(&e).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		return
	}
	found = true
	return
}

// GetEnclosuresByHash returns enclosures by hash
func GetEnclosuresByHash(hash string) (enclosures []Enclosure, err error) {
	err = DB.Where("hash = ?", hash).Find(&enclosures).Error
	return enclosures, err
}

// Delete delete enclosure e
func (e *Enclosure) Delete() error {
	// delete from storage if there is nos other occurence
	if e.StorageKey != "" {
		removeFromStore := false
		// bypass if hash ==""
		enclosures, err := GetEnclosuresByHash(e.Hash)
		if err != nil {
			return err
		}
		removeFromStore = len(enclosures) == 1
		if removeFromStore {
			if err := Store.Del(e.StorageKey); err != nil {
				return err
			}
		}
	}
	return DB.Unscoped().Delete(e).Error
}

// Update update enclosure
func (e *Enclosure) Update() error {
	return DB.Save(e).Error
}

//////////////////////////////////
// Image

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
			if !strings.HasPrefix(err.Error(), "404") {
				return err
			}
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
