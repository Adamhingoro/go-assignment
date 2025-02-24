package cronjob_test

import (
	"log"
	"os"
	"testing"
	"time"

	"assigment/cronjob"
	"assigment/database"
	"assigment/model"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// TestMain sets up the test environment
func TestMain(m *testing.M) {
	// Load environment variables from the .env file
	envErr := godotenv.Load("../test.env")
	if envErr != nil {
		log.Fatal("Error loading .env file:", envErr)
	}

	// Run the tests
	exitCode := m.Run()

	// Exit with the test result
	os.Exit(exitCode)
}

// TestCheckAndUpdateUsers tests the checkAndUpdateUsers function
func TestCheckAndUpdateUsers(t *testing.T) {
	// Arrange: Set up the test database
	config := database.Config{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
	}
	db, err := database.NewDatabase(config)
	if err != nil {
		t.Fatalf("could not connect to test database: %v", err)
	}
	defer db.DB.Exec("DELETE FROM users") // Clean up after the test

	// Auto-migrate the schema
	err = db.DB.AutoMigrate(&model.User{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Pre-populate the database with test users
	fiveHoursAgo := time.Now().Add(-6 * time.Hour).Format(time.RFC3339)
	oneHourAgo := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)

	db.DB.Exec("INSERT INTO users (email, last_seen, is_available) VALUES ($1, $2, $3)", "old@example.com", fiveHoursAgo, true)
	db.DB.Exec("INSERT INTO users (email, last_seen, is_available) VALUES ($1, $2, $3)", "recent@example.com", oneHourAgo, true)

	// Act: Call the function to update users
	cronjob.CheckAndUpdateUsers(db)

	// Assert: Verify that the old user's is_available field was updated
	var oldUser model.User
	db.DB.Where("email = ?", "old@example.com").First(&oldUser)
	assert.False(t, oldUser.IsAvailable, "Expected is_available to be false for old user")

	// Assert: Verify that the recent user's is_available field remains unchanged
	var recentUser model.User
	db.DB.Where("email = ?", "recent@example.com").First(&recentUser)
	assert.True(t, recentUser.IsAvailable, "Expected is_available to remain true for recent user")
}

// TestNewCheckAndUpdateUsersCronJob tests the initialization of the cron job
func TestNewCheckAndUpdateUsersCronJob(t *testing.T) {
	// Arrange: Set up the test database
	config := database.Config{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
	}
	db, err := database.NewDatabase(config)
	if err != nil {
		t.Fatalf("could not connect to test database: %v", err)
	}
	defer db.DB.Exec("DELETE FROM users") // Clean up after the test

	// Act: Initialize the cron job
	cronInstance := cronjob.NewCheckAndUpdateUsersCronJob(db)

	// Assert: Ensure the cron instance is not nil
	assert.NotNil(t, cronInstance, "Expected cron instance to be initialized")
}
