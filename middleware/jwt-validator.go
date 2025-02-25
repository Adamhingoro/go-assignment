package middleware

import (
	"assignment/core"
	"assignment/handler"
	"assignment/model"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func ValidateJWT(tokenString string, core *core.CoreConfig) (*jwt.Token, error) {
	log.Println("Received the jwt token", tokenString)
	// Remove the "Bearer " prefix from the token string
	tokenString = strings.Split(tokenString, " ")[1]
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return core.JWT_KEY, nil
	})
}

func ValidateJWTMiddleware(next http.HandlerFunc, core *core.CoreConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			handler.JSONError(w, http.StatusUnauthorized, "Missing token")
			return
		}

		token, err := ValidateJWT(tokenString, core)
		if err != nil || !token.Valid {
			handler.JSONError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			handler.JSONError(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		log.Println("the token claims are ", claims)

		email, ok := claims["email"].(string)
		if !ok {
			http.Error(w, "Invalid email in token", http.StatusUnauthorized)
			return
		}

		go UpdateUserLastSeenAndIsAvailable(email)

		// Restrict based on payload values
		if isAdmin, ok := claims["is_admin"].(bool); !ok || !isAdmin {
			handler.JSONError(w, http.StatusForbidden, "Access denied")
			return
		}

		next(w, r)
	}
}

func UpdateUserLastSeenAndIsAvailable(email string) {
	log.Println("Trying to update last seen of user", email)
	var user model.User
	if err := DB.First(&user, "email = ?", email).Error; err != nil {
		log.Println("Unable to find user based on email for last seen update", email)
		return
	}

	// Update the last_seen field
	if err := DB.Model(&user).Update("last_seen", time.Now()).Error; err != nil {
		log.Println("failed to update last seen of user", email)
		return
	}

	// Update the is_available field
	if err := DB.Model(&user).Update("is_available", true).Error; err != nil {
		log.Println("failed to update last seen of user", email)
		return
	}

	log.Println("Last seen updated of user ", email)
}
