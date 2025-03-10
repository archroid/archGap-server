package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email          string `gorm:"uniqueIndex;not null"`
	Password       string `gorm:"not null"`
	Username       string
	ProfilePicture string
	LastSeen       *gorm.DeletedAt 
	IsOnline       bool       
}
