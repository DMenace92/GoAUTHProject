package models

import (
	"gorm.io/gorm"
)

// User struct to define the User model

type User struct {
	gorm.Model
	Firstname string
	Lastname  string
	Email     string
	Username  string
	Password  string
	Privilege string
	Token     string
}
