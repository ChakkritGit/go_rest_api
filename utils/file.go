package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// กำหนดโฟลเดอร์หลัก
const UploadDir = "uploads"

// 1. ฟังก์ชันสำหรับ Upload ไฟล์
// คืนค่า: (path ของไฟล์ที่ save ได้, error)
func SaveFile(c *gin.Context, file *multipart.FileHeader, subFolder string) (string, error) {
	// 1. ตรวจสอบนามสกุลไฟล์ (Allow เฉพาะรูปภาพ)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true}
	if !allowExts[ext] {
		return "", errors.New("file type not allowed")
	}

	// 2. สร้างโฟลเดอร์ถ้ายังไม่มี (เช่น uploads/avatars)
	fullPath := filepath.Join(UploadDir, subFolder)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		os.MkdirAll(fullPath, 0755) // สร้างโฟลเดอร์
	}

	// 3. ตั้งชื่อไฟล์ใหม่ (ใช้ UUID หรือ Timestamp กันชื่อซ้ำ)
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// หรือใช้ Timestamp:
	// newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// 4. Save ไฟล์
	dst := filepath.Join(fullPath, newFileName)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		return "", err
	}

	// คืนค่า Path เพื่อนำไปลง Database (เช่น uploads/avatars/xxxx.jpg)
	// ปรับ Slash ให้เป็น / เสมอ (เพื่อความชัวร์ใน URL)
	return filepath.ToSlash(dst), nil
}

// 2. ฟังก์ชันลบไฟล์
func RemoveFile(filePath string) error {
	// ถ้าส่งค่าว่างมา ไม่ต้องทำอะไร
	if filePath == "" {
		return nil
	}

	// เช็คว่ามีไฟล์อยู่จริงไหม
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // ถือว่าลบไปแล้ว ไม่ error
	}

	// ลบไฟล์
	return os.Remove(filePath)
}
