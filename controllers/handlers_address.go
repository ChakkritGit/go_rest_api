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
	// 1. รับค่าจาก Form (ไม่ใช่ JSON แล้ว)
	firstname := c.PostForm("firstname")
	lastname := c.PostForm("lastname")
	codeStr := c.PostForm("code")
	codeInt, err := strconv.Atoi(codeStr)
	Phone := c.PostForm("Phone")

	var imagePath string

	// 2. รับไฟล์รูป (key ชื่อ "image")
	file, err := c.FormFile("image")
	if err == nil { // ถ้ามีการแนบไฟล์มา
		// เรียก Utils ง่ายๆ บรรทัดเดียว!
		savedPath, err := utils.SaveFile(c, file, "address")
		if err != nil {
			utils.RespondError(c, http.StatusBadRequest, err.Error())
			return
		}
		imagePath = savedPath
	}

	// 3. สร้างข้อมูล
	input := models.AddressBook{
		Firstname: firstname,
		Lastname:  lastname,
		Code:      codeInt,
		Phone:     Phone,
		Image:     imagePath, // บันทึก Path ลง DB
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

	file, err := c.FormFile("image")
	if err == nil {
		// 1. ลบรูปเก่าทิ้ง (ถ้ามี)
		utils.RemoveFile(item.Image)

		// 2. อัปโหลดรูปใหม่
		newPath, err := utils.SaveFile(c, file, "address")
		if err != nil {
			utils.RespondError(c, 400, err.Error())
			return
		}
		item.Image = newPath // อัปเดต Path
	}

	h.DB.Save(&item)

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
	id, _ := strconv.Atoi(c.Param("id"))

	var item models.AddressBook
	if err := h.DB.First(&item, id).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Not found")
		return
	}

	// [เพิ่ม] ลบรูปไฟล์จริงออกจากเครื่องก่อน
	// เรียก Utils ง่ายๆ บรรทัดเดียว!
	utils.RemoveFile(item.Image)

	h.DB.Delete(&item)
	utils.RespondJSON(c, http.StatusOK, "deleted", nil)
}
