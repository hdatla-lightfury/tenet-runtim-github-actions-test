package models

import "time"

type Account struct {
	UserID      string            `json:"user_id"`
	Username    string            `json:"username"`
	Email       string            `json:"email"`
	DisplayName string            `json:"display_name"`
	AvatarURL   string            `json:"avatar_url"`
	LangTag     string            `json:"lang_tag"`
	Location    string            `json:"location"`
	Timezone    string            `json:"timezone"`
	Metadata    map[string]string `json:"metadata"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	// Your custom fields
	Level      int               `json:"level"`
	Experience int64             `json:"experience"`
	Settings   map[string]string `json:"settings"`
}

type Wallet struct {
	Coins    int64 `json:"coins"`
	Diamonds int64 `json:"diamonds"`
}

type UpdateProfileRequest struct {
	DisplayName string `json:"display_name"`
}

type UpdateProfileResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
