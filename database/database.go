package database

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
}

type Database struct {
	DB *gorm.DB
}

var once sync.Once

func NewDatabase(config Config) (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Host, config.User, config.Password, config.DBName, config.Port,
	)

	var db *gorm.DB
	var err error

	once.Do(func() {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established")
	return &Database{DB: db}, nil
}

func (db Database) CloseDB() {
	sqlDB, err := db.DB.DB()
	if err != nil {
		log.Printf("Failed to retrieve sql.DB instance: %v", err)
		return
	}

	// Close the database connection
	if err := sqlDB.Close(); err != nil {
		log.Printf("Failed to close database connection: %v", err)
	} else {
		log.Println("Database connection closed")
	}
}
