package account

import "github.com/jinzhu/gorm"

// User is a gorm model of a user account
type User struct {
	gorm.Model
	Username     string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
}
