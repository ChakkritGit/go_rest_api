package controllers

import (
	"go_rest/models"
	"go_rest/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 1. สร้าง Struct และ Constructor
type AddressController struct {
	DB *gorm.DB
}

func NewAddressController(db *gorm.DB) *AddressController {
	return &AddressController{DB: db}
}

// 2. เปลี่ยน function เป็น method และใช้ h.DB

// GET /api/address
func (h *AddressController) GetAllAddress(c *gin.Context) {
	var list []models.AddressBook
	h.DB.Find(&list)
	utils.RespondJSON(c, http.StatusOK, "success", list)
}

// POST /api/address
func (h *AddressController) CreateAddress(c *gin.Context) {
	var input models.AddressBook
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid body")
		return
	}

	h.DB.Create(&input)
	utils.RespondJSON(c, http.StatusCreated, "created", input)
}

// GET /api/address/:id
func (h *AddressController) GetAddressByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var item models.AddressBook
	if err := h.DB.First(&item, id).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Not found")
		return
	}
	utils.RespondJSON(c, http.StatusOK, "success", item)
}

// PUT /api/address/:id
func (h *AddressController) UpdateAddress(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var item models.AddressBook
	if err := h.DB.First(&item, id).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Not found")
		return
	}

	var input models.AddressBook
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid body")
		return
	}

	item.Firstname = input.Firstname
	item.Lastname = input.Lastname
	item.Code = input.Code
	item.Phone = input.Phone

	h.DB.Save(&item)
	utils.RespondJSON(c, http.StatusOK, "updated", item)
}

// DELETE /api/address/:id
func (h *AddressController) DeleteAddress(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := h.DB.Delete(&models.AddressBook{}, id).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Not found")
		return
	}
	utils.RespondJSON(c, http.StatusOK, "deleted", nil)
}
