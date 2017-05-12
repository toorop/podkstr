package core

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// DB connector
var DB *gorm.DB

// DbAutoMigrate Auto Migrate DB (keep up2date with models)
func DbAutoMigrate() error {
	if DB == nil {
		return errors.New("DB is not initialized")
	}
	return DB.AutoMigrate(&User{}, &ShowImage{}, &Show{}, Episode{}, Enclosure{}, Image{}, Keyword{}).Error
}
