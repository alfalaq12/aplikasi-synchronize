package server

import "gorm.io/gorm"

// Struktur User untuk Database
type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string
}
