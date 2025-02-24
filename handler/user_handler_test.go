// handler/user_handler_test.go
package handler

import (
	"assigment/database"
	"assigment/model"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler(t *testing.T) {

	env_err := godotenv.Load("../test.env")
	if env_err != nil {
		log.Fatal("Error while loading env file", env_err)
	}

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

	err = db.DB.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Pre-populate the database with a test user
	testUser := model.User{Email: "test@example.com", Password: "password123"}
	db.DB.Create(&testUser)

	handler := NewUserHandler(db)

	t.Run("GetUsers", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users", nil)
		rr := httptest.NewRecorder()

		handler.GetUsers(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("CreateUser", func(t *testing.T) {
		newUser := model.User{Email: "newuser@example.com", Password: "newpassword"}
		body, _ := json.Marshal(newUser)
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.CreateUser(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		updatedUser := model.User{Email: "updated@example.com", Password: "updatedpassword"}
		body, _ := json.Marshal(updatedUser)
		req, _ := http.NewRequest("PUT", "/users?id=1", bytes.NewBuffer(body)) // Assuming ID is 1
		rr := httptest.NewRecorder()

		handler.UpdateUser(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/users?id=1", nil) // Assuming ID is 1
		rr := httptest.NewRecorder()

		handler.DeleteUser(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
	})
}
