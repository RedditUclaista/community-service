package usecases

import (
	"context"

	"github.com/RedditUclaista/community-service/internal/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CommunityRepository interface {
	Create(ctx context.Context, tx pgx.Tx, c *entities.Community) error
	Update(ctx context.Context, tx pgx.Tx, c *entities.Community) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Community, error)
	SearchPaginated(ctx context.Context, query string, limit int, offset int) ([]entities.Community, error)
	GetCommunitiesBulk(ctx context.Context, ids []uuid.UUID) ([]entities.Community, error)
}

type MemberRepository interface {
	Subscribe(ctx context.Context, tx pgx.Tx, m *entities.CommunityMember) error
	Unsubscribe(ctx context.Context, tx pgx.Tx, communityID, userID uuid.UUID) error
	ChangeRole(ctx context.Context, tx pgx.Tx, communityID, userID uuid.UUID, role entities.Role) error
	GetMember(ctx context.Context, communityID, userID uuid.UUID) (*entities.CommunityMember, error)
	GetMembers(ctx context.Context, communityID uuid.UUID) ([]entities.CommunityMember, error)
	GetUserCommunities(ctx context.Context, userID uuid.UUID) ([]entities.Community, error)
}

type OutboxRepository interface {
	Insert(ctx context.Context, tx pgx.Tx, event *entities.OutboxEvent) error
}

type UserRepository interface {
	Create(ctx context.Context, id uuid.UUID) error
}
