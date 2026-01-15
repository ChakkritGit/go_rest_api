package utils

import (
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func RespondJSON(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, APIResponse{
		Code:    status,
		Message: message,
		Data:    data,
	})
}

func RespondError(c *gin.Context, status int, message string) {
	RespondJSON(c, status, message, nil)
}
