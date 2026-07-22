package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type EventType string
type AggregateType string
type Status string

const (
	TypeCommunity       EventType = "COMMUNITY"
	TypeCommunityMember EventType = "COMMUNITY_MEMBER"
)

const (
	AggregateCommunityCreated AggregateType = "COMMUNITY_CREATED"
	AggregateCommunityUpdated AggregateType = "COMMUNITY_UPDATED"
	AggregateMemberJoined     AggregateType = "MEMBER_JOINED"
	AggregateMemberLeft       AggregateType = "MEMBER_LEFT"
)

const (
	StatusPending  Status = "PENDING"
	StatusComplete Status = "COMPLETE"
	StatusFail     Status = "FAIL"
)

type OutboxEvent struct {
	ID            uuid.UUID       `json:"id"`
	AggregateType AggregateType   `json:"aggregate_type"`
	AggregateID   uuid.UUID       `json:"aggregate_id"`
	Type          EventType       `json:"type"`
	Payload       json.RawMessage `json:"payload"`
	Status        Status          `json:"status"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}
