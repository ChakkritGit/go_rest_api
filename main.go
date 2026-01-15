package main

import (
	"go_rest/controllers"
	"go_rest/database"
	"go_rest/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. เชื่อมต่อ DB (ConnectDB จะเซ็ตค่าให้ database.DB)
	database.ConnectDB()

	// ดึง instance ของ DB มาเก็บไว้เพื่อส่งต่อ
	db := database.DB

	// 2. สร้าง Controller Instances โดยส่ง db เข้าไป (Dependency Injection)
	authCtrl := controllers.NewAuthController(db)
	addressCtrl := controllers.NewAddressController(db)

	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.Static("/uploads", "./uploads")

	// [เพิ่ม] จำกัดขนาดไฟล์ Upload ที่ 10MB (ใส่ก่อนเข้า Route อื่นๆ)
	r.Use(utils.MaxSizeMiddleware(10))

	r.NoRoute(func(c *gin.Context) {
		utils.RespondError(c, 404, "Route not found")
	})

	api := r.Group("/api")
	{
		// Auth Routes
		auth := api.Group("/auth")
		{
			// เรียกใช้ Method ผ่านตัวแปร authCtrl
			auth.POST("/register", authCtrl.Register)
			auth.POST("/login", authCtrl.Login)
		}

		// Protected Routes
		address := api.Group("/address")
		address.Use(utils.AuthMiddleware())
		{
			// เรียกใช้ Method ผ่านตัวแปร addressCtrl
			address.GET("", addressCtrl.GetAllAddress)
			address.POST("", addressCtrl.CreateAddress)
			address.GET("/:id", addressCtrl.GetAddressByID)
			address.PUT("/:id", addressCtrl.UpdateAddress)
			address.DELETE("/:id", addressCtrl.DeleteAddress)
		}
	}

	r.Run(":8080")
}
