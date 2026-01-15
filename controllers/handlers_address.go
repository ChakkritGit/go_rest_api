package controllers

import (
	"go_rest/database"
	"go_rest/models"
	"go_rest/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GET /api/address
func GetAllAddress(c *gin.Context) {
	var list []models.AddressBook
	database.DB.Find(&list)
	utils.RespondJSON(c, http.StatusOK, "success", list)
}

// POST /api/address
func CreateAddress(c *gin.Context) {
	var input models.AddressBook
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid body")
		return
	}

	database.DB.Create(&input)
	utils.RespondJSON(c, http.StatusCreated, "created", input)
}

// GET /api/address/:id
func GetAddressByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var item models.AddressBook
	if err := database.DB.First(&item, id).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Not found")
		return
	}
	utils.RespondJSON(c, http.StatusOK, "success", item)
}

// PUT /api/address/:id
func UpdateAddress(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var item models.AddressBook
	if err := database.DB.First(&item, id).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Not found")
		return
	}

	var input models.AddressBook
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid body")
		return
	}

	// อัปเดตข้อมูล
	item.Firstname = input.Firstname
	item.Lastname = input.Lastname
	item.Code = input.Code
	item.Phone = input.Phone

	database.DB.Save(&item)
	utils.RespondJSON(c, http.StatusOK, "updated", item)
}

// DELETE /api/address/:id
func DeleteAddress(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	// ใช้ DB.Delete โดยระบุ Struct ว่างและ ID
	if err := database.DB.Delete(&models.AddressBook{}, id).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Not found")
		return
	}
	utils.RespondJSON(c, http.StatusOK, "deleted", nil)
}
