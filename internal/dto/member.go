package dto

import "github.com/RedditUclaista/community-service/internal/entities"

type ChangeRoleReq struct {
	Role entities.Role `json:"role"`
}

type MemberRes struct {
	UserID   string `json:"user_id"`
	Role     string `json:"role"`
	JoinedAt string `json:"joined_at"`
}

type MemberRoleRes struct {
	Role string `json:"role"`
}
