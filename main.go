package main

import (
	"go_rest/controllers"
	"go_rest/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.ConnectDB()

	r := gin.Default() // ใช้ Default ซึ่งรวม Logger และ Recovery
	// r.SetTrustedProxies([]string{"192.168.1.2"}) // Server Proxy แยกต่างหาก (ระบุ IP ของ Proxy นั้น)
	r.SetTrustedProxies(nil)

	// Custom 404
	r.NoRoute(func(c *gin.Context) {
		utils.RespondError(c, 404, "Route not found")
	})

	api := r.Group("/api")
	{
		// Auth Routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		// Protected Routes
		address := api.Group("/address")
		address.Use(utils.AuthMiddleware()) // Middleware
		{
			address.GET("", controllers.GetAllAddress)
			address.POST("", controllers.CreateAddress)
			address.GET("/:id", controllers.GetAddressByID)
			address.PUT("/:id", controllers.UpdateAddress)
			address.DELETE("/:id", controllers.DeleteAddress)
		}
	}

	r.Run(":8080")
}
