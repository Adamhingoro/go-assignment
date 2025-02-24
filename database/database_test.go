package database

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func GetTestingDatabase(t *testing.T) *Database {
	env_err := godotenv.Load("../test.env")
	if env_err != nil {
		log.Fatal("Error while loading env file", env_err)
	}

	config := Config{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
	}

	db, err := NewDatabase(config)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	return db
}

func TestNewDatabase(t *testing.T) {
	db := GetTestingDatabase(t)

	if db == nil {
		t.Fatal("expected a database instance, got nil")
	}

	// Check if the DB field is not nil
	if db.DB == nil {
		t.Fatal("expected DB field to be initialized, got nil")
	}

	// Optionally, you can check if the DB connection is valid
	sqlDB, err := db.DB.DB()
	if err != nil {
		t.Fatalf("failed to get underlying sql.DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}
}
