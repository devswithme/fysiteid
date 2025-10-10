package domain

import (
	"time"
)

type User struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Tickets []Ticket `gorm:"foreignKey:UserID"`

	GoogleID  string  `gorm:"unique;not null"`
	AvatarURL *string `gorm:"size:255"`
	Name      string  `gorm:"size:100"`
	Email     string  `gorm:"unique;not null"`
	Username  string  `gorm:"size:50;unique"`
}

func (User) TableName() string {
	return "users"
}
