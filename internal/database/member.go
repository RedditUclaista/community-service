package database

import (
	"context"

	"github.com/RedditUclaista/community-service/internal/entities"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MemberRepository struct {
	pool *pgxpool.Pool
}

func NewMemberRepository(pool *pgxpool.Pool) *MemberRepository {
	return &MemberRepository{pool: pool}
}

func (r *MemberRepository) Subscribe(ctx context.Context, tx pgx.Tx, m *entities.CommunityMember) error {
	_, _ = tx.Exec(ctx, `INSERT INTO PUBLIC."USER" (ID) VALUES ($1) ON CONFLICT (ID) DO NOTHING;`, m.UserID)

	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto(`PUBLIC."COMMUNITY_MEMBER"`)
	ib.Cols("COMMUNITY_ID", "USER_ID", "ROLE", "JOINED_AT")
	ib.Values(m.CommunityID, m.UserID, m.Role, m.JoinedAt)
	
	sql, args := ib.Build()
	_, err := tx.Exec(ctx, sql, args...)
	return err
}

func (r *MemberRepository) Unsubscribe(ctx context.Context, tx pgx.Tx, communityID, userID uuid.UUID) error {
	db := sqlbuilder.PostgreSQL.NewDeleteBuilder()
	db.DeleteFrom(`PUBLIC."COMMUNITY_MEMBER"`)
	db.Where(
		db.Equal("COMMUNITY_ID", communityID),
		db.Equal("USER_ID", userID),
	)

	sql, args := db.Build()
	_, err := tx.Exec(ctx, sql, args...)
	return err
}

func (r *MemberRepository) ChangeRole(ctx context.Context, tx pgx.Tx, communityID, userID uuid.UUID, role entities.Role) error {
	ub := sqlbuilder.PostgreSQL.NewUpdateBuilder()
	ub.Update(`PUBLIC."COMMUNITY_MEMBER"`)
	ub.Set(
		ub.Assign("ROLE", role),
	)
	ub.Where(
		ub.Equal("COMMUNITY_ID", communityID),
		ub.Equal("USER_ID", userID),
	)

	sql, args := ub.Build()
	_, err := tx.Exec(ctx, sql, args...)
	return err
}

func (r *MemberRepository) GetMember(ctx context.Context, communityID, userID uuid.UUID) (*entities.CommunityMember, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("COMMUNITY_ID", "USER_ID", "ROLE", "JOINED_AT")
	sb.From(`PUBLIC."COMMUNITY_MEMBER"`)
	sb.Where(
		sb.Equal("COMMUNITY_ID", communityID),
		sb.Equal("USER_ID", userID),
	)

	sql, args := sb.Build()
	row := r.pool.QueryRow(ctx, sql, args...)

	var m entities.CommunityMember
	err := row.Scan(&m.CommunityID, &m.UserID, &m.Role, &m.JoinedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *MemberRepository) GetUserCommunities(ctx context.Context, userID uuid.UUID) ([]entities.Community, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("C.ID", "C.NAME", "C.DESCRIPTION", "C.RULES", "C.BANNER_URL", "C.PROFILE_URL", "C.CREATED_BY", "C.ACTIVE", "C.CREATED_AT", "CM.ROLE")
	sb.From(`PUBLIC."COMMUNITY" C`)
	sb.JoinWithOption(sqlbuilder.InnerJoin, `PUBLIC."COMMUNITY_MEMBER" CM`, "C.ID = CM.COMMUNITY_ID")
	sb.Where(sb.Equal("CM.USER_ID", userID))
	sb.OrderBy("C.CREATED_AT").Desc()

	sql, args := sb.Build()
	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var communities []entities.Community
	for rows.Next() {
		var c entities.Community
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.Rules, &c.BannerURL, &c.ProfileURL, &c.CreatedBy, &c.Active, &c.CreatedAt, &c.Role); err != nil {
			return nil, err
		}
		communities = append(communities, c)
	}
	return communities, nil
}

func (r *MemberRepository) GetMembers(ctx context.Context, communityID uuid.UUID) ([]entities.CommunityMember, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("COMMUNITY_ID", "USER_ID", "ROLE", "JOINED_AT")
	sb.From(`PUBLIC."COMMUNITY_MEMBER"`)
	sb.Where(sb.Equal("COMMUNITY_ID", communityID))

	sql, args := sb.Build()
	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []entities.CommunityMember
	for rows.Next() {
		var m entities.CommunityMember
		if err := rows.Scan(&m.CommunityID, &m.UserID, &m.Role, &m.JoinedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}
