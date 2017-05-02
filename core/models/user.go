package models

import "github.com/jinzhu/gorm"

// User represents a Podkstr user
type User struct {
	gorm.Model
	UUID      string
	FirstName string
	LastName  string
	Email     string
	Passwd    string
}

// NewUser create and return a new user
func NewUser(email, clearPasswd string) (u User, err error) {
	return
}
