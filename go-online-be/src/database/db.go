package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	Id       int    `gorm:"primaryKey"`
	Email    string `gorm:"unique;not null"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

type Session struct {
	ID     string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID int    `gorm:"not null"`
	User   User   `gorm:"constraint:OnDelete:CASCADE;"`

	Token     string    `gorm:"unique;not null"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	ExpiresAt time.Time `gorm:"not null"`
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

	if err := db.AutoMigrate(&User{}, &Session{}); err != nil {
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
