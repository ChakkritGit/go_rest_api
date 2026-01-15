package main

import (
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func respondJSON(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, APIResponse{
		Code:    status,
		Message: message,
		Data:    data,
	})
}

func respondError(c *gin.Context, status int, message string) {
	respondJSON(c, status, message, nil)
}
