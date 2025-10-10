package domain

import (
	"time"
)

type Ticket struct {
	ID        string `gorm:"primaryKey;size:6"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Title           string  `gorm:"size:64;not null"`
	Description     string  `gorm:"size:128;not null"`
	Quota           uint    `gorm:"not null"`
	RegisteredCount uint    `gorm:"not null;default:0"`
	ImageURL        *string `gorm:"size:255"`
	Mode            *bool   `gorm:"default:true"`
	UserID          uint    `gorm:"not null;index"`
	User            User    `gorm:"foreignKey:UserID"`
}

func (Ticket) TableName() string {
	return "tickets"
}
