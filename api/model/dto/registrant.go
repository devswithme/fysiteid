package dto

import "time"

type RegistrantGetByTicketID struct {
	Username   string    `json:"username"`
	AvatarURL  *string   `json:"picture"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}

type PaginatedRegistrants struct {
	Data       []RegistrantGetByTicketID `json:"data"`
	Total      int64                     `json:"total"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalPages int                       `json:"total_pages"`
}

type RegistrantGetByUserID struct {
	ID         string    `json:"id"`
	TicketID   string    `json:"ticket_id"`
	URL        string    `json:"url"`
	Title      string    `json:"title"`
	ImageURL   *string   `json:"image"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}
