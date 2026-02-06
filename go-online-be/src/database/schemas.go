package database

import (
	"time"

	"gorm.io/datatypes"
)

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

type Game struct {
	ID        string `gorm:"primaryKey;type:text"`
	BoardSize int    `gorm:"not null"`

	PlayerBlack *string `gorm:"type:text"`
	PlayerWhite *string `gorm:"type:text"`

	CurrentTurn  uint8 `gorm:"not null"` // 1=white|2=black
	Passed       bool  `gorm:"not null;default:false"`
	GameProgress uint8 `gorm:"not null;default:0"`

	WhitePoints int `gorm:"not null;default:0"`
	BlackPoints int `gorm:"not null;default:0"`

	MoveNo      int   `gorm:"not null;default:0"`
	CurrentHash int64 `gorm:"not null;default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type GameMove struct {
	GameID string `gorm:"primaryKey;type:text"`
	MoveNo int    `gorm:"primaryKey;not null"`

	Color         uint8 `gorm:"not null"` // 1=white|2=black
	X             int   // -1=pass
	Y             int
	ResultingHash int64 `gorm:"not null"`
	CreatedAt     time.Time
}

type GameSnapshot struct {
	GameID string `gorm:"primaryKey;type:text"`
	MoveNo int    `gorm:"primaryKey;not null"`

	BoardJSON   datatypes.JSON `gorm:"type:jsonb;not null"`
	CurrentTurn uint8          `gorm:"not null"`
	Passed      bool           `gorm:"not null"`
	Hash        int64          `gorm:"not null"`

	CreatedAt time.Time
}
