package config

import (
	"fmt"
	"log"

	"incident-tracker-api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=127.0.0.1 user=postgres password=postgres dbname=incident_tracker port=5432 sslmode=disable"

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal(err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("[GORM] Connected to PostgreSQL successfully")

	err = database.AutoMigrate(
		&models.Incident{},
		&models.User{},
	)

	if err != nil {
		log.Fatal("Auto migration failed:", err)
	}

	fmt.Println("[GORM] Auto-migration completed")

	DB = database
}
