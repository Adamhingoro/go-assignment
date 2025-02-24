// models/user.go
package model

import (
	"assigment/helper"
	"log"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name        string
	Email       string `gorm:"unique"`
	Password    string
	CompanyID   int
	IsAdmin     bool
	IsAvailable bool
	LastSeen    time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if u.Password != "" {
		log.Printf("Auto hashing the password for the user")
		hashedPassword := helper.HashString(u.Password)
		u.Password = string(hashedPassword)
	}
	return nil
}
