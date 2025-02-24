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

func TestCompanyHandler(t *testing.T) {
	// Load environment variables from a test-specific .env file
	envErr := godotenv.Load("../test.env")
	if envErr != nil {
		log.Fatal("Error while loading env file", envErr)
	}

	// Database configuration
	config := database.Config{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
	}

	// Connect to the test database
	db, err := database.NewDatabase(config)
	if err != nil {
		t.Fatalf("could not connect to test database: %v", err)
	}

	// Auto-migrate the Company model
	err = db.DB.AutoMigrate(&model.Company{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Pre-populate the database with a test company
	testCompany := model.Company{Name: "Test Company", Address: "123 Test St"}
	db.DB.Create(&testCompany)

	// Initialize the CompanyHandler
	handler := NewCompanyHandler(db)

	// Test cases
	t.Run("GetCompanies", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/companies", nil)
		rr := httptest.NewRecorder()
		handler.GetCompanies(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		// Optionally, validate the response body
		var companies []model.Company
		err := json.Unmarshal(rr.Body.Bytes(), &companies)
		assert.NoError(t, err)
		assert.NotEmpty(t, companies)
	})

	t.Run("CreateCompany", func(t *testing.T) {
		newCompany := model.Company{Name: "New Company", Address: "456 New St"}
		body, _ := json.Marshal(newCompany)

		req, _ := http.NewRequest("POST", "/companies", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		handler.CreateCompany(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		// Optionally, validate the response body
		var createdCompany model.Company
		err := json.Unmarshal(rr.Body.Bytes(), &createdCompany)
		assert.NoError(t, err)
		assert.Equal(t, newCompany.Name, createdCompany.Name)
		assert.Equal(t, newCompany.Address, createdCompany.Address)
	})

	t.Run("UpdateCompany", func(t *testing.T) {
		updatedCompany := model.Company{Name: "Updated Company", Address: "789 Updated St"}
		body, _ := json.Marshal(updatedCompany)

		// Assuming the ID of the pre-populated test company is 1
		req, _ := http.NewRequest("PUT", "/companies?id=1", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		handler.UpdateCompany(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		// Optionally, validate the response body
		var updatedResponse model.Company
		err := json.Unmarshal(rr.Body.Bytes(), &updatedResponse)
		assert.NoError(t, err)
		assert.Equal(t, updatedCompany.Name, updatedResponse.Name)
		assert.Equal(t, updatedCompany.Address, updatedResponse.Address)
	})

	t.Run("DeleteCompany", func(t *testing.T) {
		// Assuming the ID of the pre-populated test company is 1
		req, _ := http.NewRequest("DELETE", "/companies?id=1", nil)
		rr := httptest.NewRecorder()
		handler.DeleteCompany(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)

		// Optionally, verify that the company was deleted
		var company model.Company
		result := db.DB.Where("id = ?", 1).First(&company)
		assert.Error(t, result.Error) // Should return an error because the record no longer exists
	})
}
