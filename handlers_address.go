package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GET /api/address
func getAllAddress(c *gin.Context) {
	var list []AddressBook
	DB.Find(&list)
	respondJSON(c, http.StatusOK, "success", list)
}

// POST /api/address
func createAddress(c *gin.Context) {
	var input AddressBook
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, http.StatusBadRequest, "Invalid body")
		return
	}

	DB.Create(&input)
	respondJSON(c, http.StatusCreated, "created", input)
}

// GET /api/address/:id
func getAddressByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var item AddressBook
	if err := DB.First(&item, id).Error; err != nil {
		respondError(c, http.StatusNotFound, "Not found")
		return
	}
	respondJSON(c, http.StatusOK, "success", item)
}

// PUT /api/address/:id
func updateAddress(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var item AddressBook
	if err := DB.First(&item, id).Error; err != nil {
		respondError(c, http.StatusNotFound, "Not found")
		return
	}

	var input AddressBook
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, http.StatusBadRequest, "Invalid body")
		return
	}

	// อัปเดตข้อมูล
	item.Firstname = input.Firstname
	item.Lastname = input.Lastname
	item.Code = input.Code
	item.Phone = input.Phone

	DB.Save(&item)
	respondJSON(c, http.StatusOK, "updated", item)
}

// DELETE /api/address/:id
func deleteAddress(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	// ใช้ DB.Delete โดยระบุ Struct ว่างและ ID
	if err := DB.Delete(&AddressBook{}, id).Error; err != nil {
		respondError(c, http.StatusNotFound, "Not found")
		return
	}
	respondJSON(c, http.StatusOK, "deleted", nil)
}
