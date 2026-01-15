package models

import "gorm.io/gorm"

// User สำหรับ Auth
type User struct {
	gorm.Model
	Username string `gorm:"unique" json:"username"`
	Password string `json:"-"`
}

// AddressBook
type AddressBook struct {
	gorm.Model
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Code      int    `json:"code"`
	Phone     string `json:"phone"`
	Image     string `json:"image"`
}
