package controllers

import (
	"go_rest/models"
	"go_rest/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AddressController struct {
	DB *gorm.DB
}

func NewAddressController(db *gorm.DB) *AddressController {
	return &AddressController{DB: db}
}

// GetAllAddress godoc
// @Summary      Get all addresses
// @Description  Retrieve a list of all addresses
// @Tags         Address
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.AddressBook
// @Failure      500  {object}  utils.APIResponse
// @Router       /address [get]
func (h *AddressController) GetAllAddress(c *gin.Context) {
	var list []models.AddressBook
	h.DB.Find(&list)
	utils.RespondJSON(c, http.StatusOK, "success", list)
}

// CreateAddress godoc
// @Summary      Create a new address
// @Description  Create address with image upload
// @Tags         Address
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        firstname formData string true "Firstname"
// @Param        lastname  formData string true "Lastname"
// @Param        code      formData int    true "Code"
// @Param        phone     formData string true "Phone"
// @Param        image     formData file   false "Address Image"
// @Success      201  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Router       /address [post]
func (h *AddressController) CreateAddress(c *gin.Context) {
	firstname := c.PostForm("firstname")
	lastname := c.PostForm("lastname")
	phone := c.PostForm("phone")

	codeStr := c.PostForm("code")
	codeInt, _ := strconv.Atoi(codeStr)

	var imagePath string

	file, err := c.FormFile("image")
	if err == nil {
		savedPath, err := utils.SaveFile(c, file, "address")
		if err != nil {
			utils.RespondError(c, http.StatusBadRequest, err.Error())
			return
		}
		imagePath = savedPath
	}

	input := models.AddressBook{
		Firstname: firstname,
		Lastname:  lastname,
		Code:      codeInt,
		Phone:     phone,
		Image:     imagePath,
	}

	h.DB.Create(&input)
	utils.RespondJSON(c, http.StatusCreated, "created", input)
}

// GetAddressByID godoc
// @Summary      Get address by ID
// @Description  Retrieve address by ID
// @Tags         Address
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Address ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /address/{id} [get]
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

// UpdateAddress godoc
// @Summary      Update address by ID
// @Description  Update address with optional image upload
// @Tags         Address
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        id        path     int    true  "Address ID"
// @Param        firstname formData string true  "Firstname"
// @Param        lastname  formData string true  "Lastname"
// @Param        code      formData int    true  "Code"
// @Param        phone     formData string true  "Phone"
// @Param        image     formData file   false "Address Image"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /address/{id} [put]
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

	file, err := c.FormFile("image")
	if err == nil {
		utils.RemoveFile(item.Image)

		newPath, err := utils.SaveFile(c, file, "address")
		if err != nil {
			utils.RespondError(c, http.StatusBadRequest, err.Error())
			return
		}
		item.Image = newPath
	}

	item.Firstname = c.PostForm("firstname")
	item.Lastname = c.PostForm("lastname")
	item.Phone = c.PostForm("phone")

	if codeStr := c.PostForm("code"); codeStr != "" {
		codeInt, _ := strconv.Atoi(codeStr)
		item.Code = codeInt
	}

	h.DB.Save(&item)
	utils.RespondJSON(c, http.StatusOK, "updated", item)
}

// DeleteAddress godoc
// @Summary      Delete address by ID
// @Description  Delete address and associated image
// @Tags         Address
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Address ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /address/{id} [delete]
func (h *AddressController) DeleteAddress(c *gin.Context) {
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

	utils.RemoveFile(item.Image)

	h.DB.Delete(&item)
	utils.RespondJSON(c, http.StatusOK, "deleted", nil)
}
