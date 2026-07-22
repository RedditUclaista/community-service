package entities

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleMember    Role = "MEMBER"
	RoleModerator Role = "MODERATOR"
)

type CommunityMember struct {
	CommunityID uuid.UUID `json:"community_id"`
	UserID      uuid.UUID `json:"user_id"`
	Role        Role      `json:"role"`
	JoinedAt    time.Time `json:"joined_at"`
}
