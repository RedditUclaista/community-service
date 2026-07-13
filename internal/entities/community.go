package entities

import (
	"time"
	"github.com/google/uuid"
)

type Community struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by"` // user email or id
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
}
