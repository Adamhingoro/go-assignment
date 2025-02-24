package middleware

import (
	"log"

	"gorm.io/gorm"
)

var DB *gorm.DB

// SetDB allows you to set the DB instance for the middleware.
func SetDB(db *gorm.DB) {
	log.Println("Setting DB for middleware")
	DB = db
}
