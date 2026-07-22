package usecases

import (
	"context"
	"encoding/json"
	"time"

	"github.com/RedditUclaista/community-service/internal/dto"
	"github.com/RedditUclaista/community-service/internal/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommunityUseCase struct {
	pool       *pgxpool.Pool
	commRepo   CommunityRepository
	outboxRepo OutboxRepository
	memberRepo MemberRepository
}

func NewCommunityUseCase(pool *pgxpool.Pool, commRepo CommunityRepository, outboxRepo OutboxRepository, memberRepo MemberRepository) *CommunityUseCase {
	return &CommunityUseCase{
		pool:       pool,
		commRepo:   commRepo,
		outboxRepo: outboxRepo,
		memberRepo: memberRepo,
	}
}

func (uc *CommunityUseCase) Create(ctx context.Context, req dto.CreateCommunityReq, userID uuid.UUID) (*entities.Community, error) {
	tx, err := uc.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	comm := &entities.Community{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Rules:       req.Rules,
		BannerURL:   req.BannerURL,
		ProfileURL:  req.ProfileURL,
		CreatedBy:   userID,
		Active:      true,
		CreatedAt:   time.Now().UTC(),
	}

	if err := uc.commRepo.Create(ctx, tx, comm); err != nil {
		return nil, err
	}

	payload, _ := json.Marshal(comm)
	event := &entities.OutboxEvent{
		ID:            uuid.New(),
		AggregateType: entities.AggregateCommunityCreated,
		AggregateID:   comm.ID,
		Type:          entities.TypeCommunity,
		Payload:       payload,
		Status:        entities.StatusPending,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	if err := uc.outboxRepo.Insert(ctx, tx, event); err != nil {
		return nil, err
	}

	member := &entities.CommunityMember{
		CommunityID: comm.ID,
		UserID:      userID,
		Role:        entities.RoleModerator,
		JoinedAt:    time.Now().UTC(),
	}

	if err := uc.memberRepo.Subscribe(ctx, tx, member); err != nil {
		return nil, err
	}

	memberPayload, _ := json.Marshal(member)
	memberEvent := &entities.OutboxEvent{
		ID:            uuid.New(),
		AggregateType: entities.AggregateMemberJoined,
		AggregateID:   comm.ID,
		Type:          entities.TypeCommunityMember,
		Payload:       memberPayload,
		Status:        entities.StatusPending,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	if err := uc.outboxRepo.Insert(ctx, tx, memberEvent); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return comm, nil
}

func (uc *CommunityUseCase) Update(ctx context.Context, id uuid.UUID, req dto.UpdateCommunityReq) (*entities.Community, error) {
	tx, err := uc.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	comm, err := uc.commRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if comm == nil {
		return nil, nil // Not found
	}

	if req.Description != nil {
		comm.Description = *req.Description
	}
	if req.Rules != nil {
		comm.Rules = *req.Rules
	}
	if req.BannerURL != nil {
		comm.BannerURL = *req.BannerURL
	}
	if req.ProfileURL != nil {
		comm.ProfileURL = *req.ProfileURL
	}

	if err := uc.commRepo.Update(ctx, tx, comm); err != nil {
		return nil, err
	}

	payload, _ := json.Marshal(comm)
	event := &entities.OutboxEvent{
		ID:            uuid.New(),
		AggregateType: entities.AggregateCommunityUpdated,
		AggregateID:   comm.ID,
		Type:          entities.TypeCommunity,
		Payload:       payload,
		Status:        entities.StatusPending,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	if err := uc.outboxRepo.Insert(ctx, tx, event); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return comm, nil
}

func (uc *CommunityUseCase) SearchPaginated(ctx context.Context, query string, limit, offset int) ([]entities.Community, error) {
	return uc.commRepo.SearchPaginated(ctx, query, limit, offset)
}
