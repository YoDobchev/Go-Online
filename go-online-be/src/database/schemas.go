package database

import "time"

type User struct {
	ID       int    `gorm:"primaryKey"`
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
