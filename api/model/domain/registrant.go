package domain

import (
	"time"
)

type Registrant struct {
	ID        string `gorm:"primaryKey;size:6"`
	CreatedAt time.Time
	UpdatedAt time.Time

	UserID     uint   `gorm:"not null;index:idx_user_ticket,unique"`
	User       User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	TicketID   string `gorm:"not null;index:idx_user_ticket,unique"`
	Ticket     Ticket `gorm:"foreignKey:TicketID;constraint:OnDelete:CASCADE"`
	IsVerified bool   `gorm:"default:false"`

	TicketOwnerID uint `gorm:"not null;index"`
}

func (Registrant) TableName() string {
	return "registrants"
}
