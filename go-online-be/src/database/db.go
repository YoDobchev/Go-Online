package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	Username string `gorm:"primaryKey" json:"username"`
	Password string `json:"password"`
}

func Connect() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env: %v", err)
	}

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
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error loading db: %v", err)
	}

	DB = db

	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatalf("auto-migrate failed: %v", err)
	}
	ctx := context.Background()

	var user User
	if err := db.WithContext(ctx).Take(&user).Error; err != nil {
		log.Println("could not get user:", err)
	} else {
		fmt.Println("user from db:", user)
	}
}
