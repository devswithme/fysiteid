package dto

type TicketCreate struct {
	Title       string `form:"title" validate:"required,min=16,max=64"`
	Description string `form:"description" validate:"required,min=32,max=128"`
	Mode        bool   `form:"mode"`
	Quota       uint   `form:"quota" validate:"required"`
}

type TicketGet struct {
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	Mode            bool    `json:"mode"`
	RegisteredCount uint    `json:"registered_count"`
	Quota           uint    `json:"quota"`
	Image           *string `json:"image"`
}

type TicketPublicGet struct {
	ID    string  `json:"id"`
	Image *string `json:"image"`
	Title string  `json:"title"`
}
