package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// POST /api/auth/register
func register(w http.ResponseWriter, r *http.Request) {
	var input User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	user := User{
		Username: input.Username,
		Password: string(hash),
	}

	if err := DB.Create(&user).Error; err != nil {
		respondError(w, http.StatusBadRequest, "Username already exists")
		return
	}

	respondJSON(w, http.StatusCreated, "registered", nil)
}

// POST /api/auth/login
func login(w http.ResponseWriter, r *http.Request) {
	var input User
	var user User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Create JWT
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString(jwtSecret)

	respondJSON(w, http.StatusOK, "success", map[string]string{
		"token": tokenStr,
	})
}
