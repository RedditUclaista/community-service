package dto

import "github.com/google/uuid"

type CreateCommunityReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Rules       string `json:"rules"`
	BannerURL   string `json:"banner_url"`
	ProfileURL  string `json:"profile_url"`
}

type UpdateCommunityReq struct {
	Description *string `json:"description"`
	Rules       *string `json:"rules"`
	BannerURL   *string `json:"banner_url"`
	ProfileURL  *string `json:"profile_url"`
}

type CommunityListRes struct {
	Communities []CommunityRes `json:"communities"`
	NextCursor  *string        `json:"next_cursor,omitempty"`
}

type CommunityRes struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Rules       string    `json:"rules"`
	BannerURL   string    `json:"banner_url"`
	ProfileURL  string    `json:"profile_url"`
	CreatedBy   uuid.UUID `json:"created_by"`
	Role        string    `json:"role,omitempty"`
}
