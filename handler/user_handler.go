package handler

import (
	"assigment/database"
	"assigment/model"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type UserHandler struct {
	DB *database.Database
}

func NewUserHandler(db *database.Database) *UserHandler {
	return &UserHandler{DB: db}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []model.User
	h.DB.DB.Find(&users)
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	json.NewDecoder(r.Body).Decode(&user)

	if user.IsAdmin {
		user.CompanyID = 0
	} else if user.CompanyID == 0 {
		JSONError(w, http.StatusBadRequest, "Invalid Company ID")
		return
	}

	userExists, _ := h.UserExistsByEmail(user.Email)
	if userExists {
		JSONError(w, http.StatusBadRequest, "User already exists with same email")

		return
	}

	h.DB.DB.Create(&user)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UserExistsByEmail(email string) (bool, error) {
	var user model.User
	result := h.DB.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil // User does not exist
		}
		return false, result.Error // Some other error occurred
	}
	return true, nil // User exists
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	json.NewDecoder(r.Body).Decode(&user)

	// Assuming the user ID is passed in the URL
	id := r.URL.Query().Get("id")
	if id == "" {
		JSONError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	result := h.DB.DB.Model(&user).Where("id = ?", id).Updates(user)
	if result.Error != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Assuming the user ID is passed in the URL
	id := r.URL.Query().Get("id")
	if id == "" {
		JSONError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	result := h.DB.DB.Where("id = ? AND is_admin = ?", id, false).Delete(&model.User{})
	if result.Error != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	w.WriteHeader(http.StatusNoContent) // No content to return
}
