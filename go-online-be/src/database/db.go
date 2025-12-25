package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf(
		"host=localhost user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Europe/Sofia",
		dbUser,
		dbPass,
		dbName,
		dbPort,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error loading db: %v", err)
	}

	if err := DB.AutoMigrate(&User{}, &Session{}); err != nil {
		log.Fatalf("auto-migrate failed: %v", err)
	}
}
