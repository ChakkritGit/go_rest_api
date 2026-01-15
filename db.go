package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func connectDB() {
	dsn := "host=localhost user= password= dbname=mydb port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// Auto Migrate
	database.AutoMigrate(&User{}, &AddressBook{})

	DB = database
	log.Println("Database connected")
}
