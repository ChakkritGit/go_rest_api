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

const UploadDir = "uploads"

func SaveFile(c *gin.Context, file *multipart.FileHeader, subFolder string) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true}
	if !allowExts[ext] {
		return "", errors.New("file type not allowed")
	}

	fullPath := filepath.Join(UploadDir, subFolder)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		os.MkdirAll(fullPath, 0755)
	}

	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	dst := filepath.Join(fullPath, newFileName)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		return "", err
	}

	return filepath.ToSlash(dst), nil
}

func RemoveFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(filePath)
}
