package handler

import (
	"assigment/core"
	"assigment/database"
	"assigment/dto"
	"assigment/helper"
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

func TestLogin(t *testing.T) {

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

	// Pre-populate the database with a test user for the successful login case
	db.DB.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", "test@example.com", helper.HashString("password123"))
	defer db.DB.Exec("DELETE FROM users") // Clean up after the test

	coreConfig := core.CoreConfig{
		JWT_KEY: []byte(os.Getenv("JWT_KEY")),
	}
	handler := NewAuthHandler(db, &coreConfig)

	tests := []struct {
		name           string
		creds          dto.LoginRequest
		expectedStatus int
	}{
		{
			name: "Successful Login",
			creds: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid Credentials",
			creds: dto.LoginRequest{
				Email:    "wrong@example.com",
				Password: "wrongpassword",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.creds)
			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			handler.Login(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}

}
