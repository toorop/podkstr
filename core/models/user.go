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
