package dto

type UserCreate struct {
	GoogleID  string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"picture"`
}

type UserGet struct {
	Name      string  `json:"name"`
	Username  string  `json:"username"`
	AvatarURL *string `json:"picture"`
}

type UserUpdate struct {
	Name     string `form:"name" validate:"required,min=1,max=36"`
	Username string `form:"username" validate:"required,min=1,max=36,slug"`
}
