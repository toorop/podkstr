package core

import (
	"errors"
	"fmt"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"

	"golang.org/x/crypto/bcrypt"
)

// User represents a Podkstr user
type User struct {
	gorm.Model
	UUID           string `gorm:"type:char(36);unique_index"`
	FirstName      string
	LastName       string
	Email          string `gorm:"unique_index"`
	Passwd         string
	ValidationUUID string `gorm:"type:char(36);unique_index"`
	Activated      bool
	Shows          []Show
}

// UserNew create and return a new user
func UserNew(email, clearPasswd string) (u User, err error) {
	if email == "" || clearPasswd == "" {
		return u, errors.New("core.UserNew - email or passwd or both missing")
	}
	if !govalidator.IsEmail(email) {
		return u, fmt.Errorf("core.UserNew - %s  is not a valid email", email)
	}
	email = strings.ToLower(email)

	// chriffrement du mail
	passwd, err := bcrypt.GenerateFromPassword([]byte(clearPasswd), 10)
	if err != nil {
		return u, err
	}
	u.Email = email
	u.Passwd = string(passwd)
	u.ValidationUUID = uuid.NewV4().String()
	u.UUID = uuid.NewV4().String()
	err = DB.Create(&u).Error
	return u, err
}

// UserGetByMail get user by mail (if exists)
func UserGetByMail(email string) (u User, found bool, err error) {
	email = strings.ToLower(strings.TrimSpace(email))
	err = DB.Where("email = ?", email).First(&u).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		return
	}
	found = true
	return
}

// UserGetByEmailPasswd return user by email and password (if exists)
func UserGetByEmailPasswd(email, passwd string) (u User, found bool, err error) {
	u, found, err = UserGetByMail(email)
	if err != nil || !found {
		return
	}
	// check passwd
	err = bcrypt.CompareHashAndPassword([]byte(u.Passwd), []byte(passwd))
	if err != nil {
		return
	}
	// User exists
	found = true
	return
}

// UserGetByValidationUUID get user by its validation UUID
func UserGetByValidationUUID(uuid string) (u User, found bool, err error) {
	err = DB.Where("validation_uuid = ?", uuid).First(&u).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		return
	}
	found = true
	return
}

// SetPasswd set user password
func (u *User) SetPasswd(clearPasswd string) error {
	passwd, err := bcrypt.GenerateFromPassword([]byte(clearPasswd), 10)
	if err != nil {
		return err
	}
	u.Passwd = string(passwd)
	return u.Save()
}

// GetShows returns User shows
func (u User) GetShows() (shows []Show, err error) {
	err = DB.Model(&u).Related(&shows).Error
	return
}

// GetShowByFeed returns a show by is feed URL
func (u User) GetShowByFeed(url string) (show Show, found bool, err error) {
	url = strings.ToLower(strings.TrimSpace(url))
	err = DB.Model(&u).Related(&Show{}).Where("feed = ?", url).First(&show).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		return
	}
	found = true
	return
}

// GetShowByUUID return an user show by his uuid
func (u User) GetShowByUUID(uuid string) (show Show, found bool, err error) {
	err = DB.Model(&u).Related(&Show{}).Where("uuid = ?", uuid).First(&show).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		return
	}
	found = true
	return
}

// ResetValidationUUID generate a new validation UUID
func (u *User) ResetValidationUUID() error {
	u.ValidationUUID = uuid.NewV4().String()
	return u.Save()
}

// Save saves user
func (u *User) Save() error {
	return DB.Save(u).Error
}
