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

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

// Register godoc
// @Summary      Register a new user
// @Description  Register a new user with username and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input body models.LoginAndRegisterInput true "User Credentials"
// @Success      201  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Router       /auth/register [post]
func (h *AuthController) Register(c *gin.Context) {
	var input models.LoginAndRegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	user := models.User{
		Username: input.Username,
		Password: string(hash),
	}

	if err := h.DB.Create(&user).Error; err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Username already exists")
		return
	}

	utils.RespondJSON(c, http.StatusCreated, "registered", nil)
}

// Login godoc
// @Summary      Login user
// @Description  Login to get JWT Token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input body models.LoginAndRegisterInput true "User Credentials"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse "Invalid Input"
// @Failure      401  {object}  utils.APIResponse "Invalid username or password"
// @Router       /auth/login [post]
func (h *AuthController) Login(c *gin.Context) {
	var input models.LoginAndRegisterInput
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := h.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Invalid username or password")
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
