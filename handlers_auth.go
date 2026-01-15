package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// POST /api/auth/register
func register(c *gin.Context) {
	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	user := User{
		Username: input.Username,
		Password: string(hash),
	}

	if err := DB.Create(&user).Error; err != nil {
		respondError(c, http.StatusBadRequest, "Username already exists")
		return
	}

	respondJSON(c, http.StatusCreated, "registered", nil)
}

// POST /api/auth/login
func login(c *gin.Context) {
	var input User
	var user User

	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		respondError(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		respondError(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Create JWT
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString(jwtSecret)

	respondJSON(c, http.StatusOK, "success", map[string]string{
		"token": tokenStr,
	})
}
