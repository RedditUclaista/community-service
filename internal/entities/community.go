package entities

import (
	"time"

	"github.com/google/uuid"
)

type Community struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Rules       string    `json:"rules"`
	BannerURL   string    `json:"banner_url"`
	ProfileURL  string    `json:"profile_url"`
	CreatedBy   uuid.UUID `json:"created_by"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	Role        string    `json:"role,omitempty"`
}
