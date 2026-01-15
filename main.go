package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	connectDB()

	r := gin.Default() // ใช้ Default ซึ่งรวม Logger และ Recovery

	// Custom 404
	r.NoRoute(func(c *gin.Context) {
		respondError(c, 404, "Route not found")
	})

	api := r.Group("/api")
	{
		// Auth Routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", register)
			auth.POST("/login", login)
		}

		// Protected Routes
		address := api.Group("/address")
		address.Use(AuthMiddleware()) // Middleware
		{
			address.GET("", getAllAddress)
			address.POST("", createAddress)
			address.GET("/:id", getAddressByID)
			address.PUT("/:id", updateAddress)
			address.DELETE("/:id", deleteAddress)
		}
	}

	log.Println("Server running at :8080")
	r.Run(":8080")
}
