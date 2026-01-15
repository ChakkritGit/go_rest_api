package database

import (
	"fmt"
	"log"
	"os"

	"go_rest/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found (loading from system env or defaulting)")
	}

	host := os.Getenv("DATABASE_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DATABASE_PORT")
	if port == "" {
		port = "5432"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=mydb port=%s sslmode=disable",
		host,
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		port,
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	database.AutoMigrate(&models.User{}, &models.AddressBook{})

	DB = database
	log.Println("Database connected")
}
