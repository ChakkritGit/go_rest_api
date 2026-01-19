package main

import (
	"go_rest/controllers"
	"go_rest/database"
	"go_rest/utils"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "go_rest/docs"
)

// @title           Go REST API
// @version         1.0
// @description     This is a sample server for managing addresses.
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:8080
// @BasePath        /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	database.ConnectDB()

	db := database.DB

	authCtrl := controllers.NewAuthController(db)
	addressCtrl := controllers.NewAddressController(db)

	docs.SwaggerInfo.Title = "My Go REST API"
	docs.SwaggerInfo.Description = "This is a sample server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api"

	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.Use(utils.MaxSizeMiddleware(10))

	r.NoRoute(func(c *gin.Context) {
		utils.RespondError(c, 404, "Route not found")
	})

	api := r.Group("/api")

	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api.Static("/uploads", "./uploads")
	{
		api.GET("/healthcheck", controllers.HealthCheckHandler)

		// Auth Routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authCtrl.Register)
			auth.POST("/login", authCtrl.Login)
		}

		// Protected Routes
		address := api.Group("/address")
		address.Use(utils.AuthMiddleware())
		{
			address.GET("", addressCtrl.GetAllAddress)
			address.POST("", addressCtrl.CreateAddress)
			address.GET("/:id", addressCtrl.GetAddressByID)
			address.PUT("/:id", addressCtrl.UpdateAddress)
			address.DELETE("/:id", addressCtrl.DeleteAddress)
		}
	}

	log.Fatal((r.Run(":8080")))
}
