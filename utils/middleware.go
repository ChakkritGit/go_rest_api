package utils

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GetJWTSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

// JWT Middleware for Gin
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			RespondError(c, 401, "Missing token")
			c.Abort() // หยุดการทำงานของ Handler ถัดไป
			return
		}

		// Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			RespondError(c, 401, "Invalid token format")
			c.Abort()
			return
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return GetJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			RespondError(c, 401, "Invalid token")
			c.Abort()
			return
		}

		// ถ้าต้องการเก็บ user_id ไว้ใช้ต่อ
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
		}

		c.Next() // ไปยัง Handler ถัดไป
	}
}

func MaxSizeMiddleware(limitMB int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limitMB*1024*1024)
		c.Next()
	}
}
