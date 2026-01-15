package controllers

import (
	"go_rest/models"
	"go_rest/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 1. สร้าง Struct และ Constructor
type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

// 2. เปลี่ยน function เป็น method (มี h *AuthController นำหน้า)
// และใช้ h.DB แทน database.DB

// POST /api/auth/register
func (h *AuthController) Register(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	user := models.User{
		Username: input.Username,
		Password: string(hash),
	}

	// ใช้ h.DB
	if err := h.DB.Create(&user).Error; err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Username already exists")
		return
	}

	utils.RespondJSON(c, http.StatusCreated, "registered", nil)
}

// POST /api/auth/login
func (h *AuthController) Login(c *gin.Context) {
	var input models.User
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	// ใช้ h.DB
	if err := h.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString(utils.GetJWTSecret())

	utils.RespondJSON(c, http.StatusOK, "success", map[string]string{
		"token": tokenStr,
	})
}
