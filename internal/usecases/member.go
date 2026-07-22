package usecases

import (
	"context"
	"encoding/json"
	"time"

	"github.com/RedditUclaista/community-service/internal/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MemberUseCase struct {
	pool       *pgxpool.Pool
	memberRepo MemberRepository
	outboxRepo OutboxRepository
}

func NewMemberUseCase(pool *pgxpool.Pool, memberRepo MemberRepository, outboxRepo OutboxRepository) *MemberUseCase {
	return &MemberUseCase{
		pool:       pool,
		memberRepo: memberRepo,
		outboxRepo: outboxRepo,
	}
}

func (uc *MemberUseCase) Join(ctx context.Context, communityID, userID uuid.UUID, role entities.Role) error {
	tx, err := uc.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	member := &entities.CommunityMember{
		CommunityID: communityID,
		UserID:      userID,
		Role:        role,
		JoinedAt:    time.Now().UTC(),
	}

	if err := uc.memberRepo.Subscribe(ctx, tx, member); err != nil {
		return err // Could be 409 Conflict if already exists
	}

	payload, _ := json.Marshal(member)
	event := &entities.OutboxEvent{
		ID:            uuid.New(),
		AggregateType: entities.AggregateMemberJoined,
		AggregateID:   communityID,
		Type:          entities.TypeCommunityMember,
		Payload:       payload,
		Status:        entities.StatusPending,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	if err := uc.outboxRepo.Insert(ctx, tx, event); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (uc *MemberUseCase) Leave(ctx context.Context, communityID, userID uuid.UUID) error {
	tx, err := uc.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := uc.memberRepo.Unsubscribe(ctx, tx, communityID, userID); err != nil {
		return err
	}

	payload, _ := json.Marshal(map[string]string{
		"community_id": communityID.String(),
		"user_id":      userID.String(),
	})
	
	event := &entities.OutboxEvent{
		ID:            uuid.New(),
		AggregateType: entities.AggregateMemberLeft,
		AggregateID:   communityID,
		Type:          entities.TypeCommunityMember,
		Payload:       payload,
		Status:        entities.StatusPending,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	if err := uc.outboxRepo.Insert(ctx, tx, event); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (uc *MemberUseCase) ChangeRole(ctx context.Context, communityID, userID uuid.UUID, role entities.Role) error {
	tx, err := uc.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := uc.memberRepo.ChangeRole(ctx, tx, communityID, userID, role); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (uc *MemberUseCase) GetUserCommunities(ctx context.Context, userID uuid.UUID) ([]entities.Community, error) {
	return uc.memberRepo.GetUserCommunities(ctx, userID)
}

func (uc *MemberUseCase) GetMembers(ctx context.Context, communityID uuid.UUID) ([]entities.CommunityMember, error) {
	return uc.memberRepo.GetMembers(ctx, communityID)
}

func (uc *MemberUseCase) GetMemberRole(ctx context.Context, communityID, userID uuid.UUID) (*entities.CommunityMember, error) {
	return uc.memberRepo.GetMember(ctx, communityID, userID)
}
