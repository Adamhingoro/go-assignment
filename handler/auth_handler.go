package handler

import (
	"assigment/core"
	"assigment/database"
	"assigment/dto"
	"assigment/helper"
	"assigment/model"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey []byte

type AuthHandler struct {
	DB *database.Database
}

func NewAuthHandler(db *database.Database, core *core.CoreConfig) *AuthHandler {
	if core.JWT_KEY == nil {
		log.Panicln("JWT key not supplied")
	}
	jwtKey = core.JWT_KEY
	return &AuthHandler{DB: db}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var user model.User
	if err := h.DB.DB.Where("email = ? AND password = ?", creds.Email, helper.HashString(creds.Password)).First(&user).Error; err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid Email or Password")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    user.Email,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Error while generating JWT token")
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
