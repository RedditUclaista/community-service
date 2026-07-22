package database

import (
	"context"
	"strings"

	"github.com/RedditUclaista/community-service/internal/entities"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommunityRepository struct {
	pool *pgxpool.Pool
}

func NewCommunityRepository(pool *pgxpool.Pool) *CommunityRepository {
	return &CommunityRepository{pool: pool}
}

func (r *CommunityRepository) Create(ctx context.Context, tx pgx.Tx, c *entities.Community) error {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto(`PUBLIC."COMMUNITY"`)
	ib.Cols("ID", "NAME", "DESCRIPTION", "RULES", "BANNER_URL", "PROFILE_URL", "CREATED_BY", "ACTIVE", "CREATED_AT")
	ib.Values(c.ID, c.Name, c.Description, c.Rules, c.BannerURL, c.ProfileURL, c.CreatedBy, c.Active, c.CreatedAt)

	sql, args := ib.Build()
	_, err := tx.Exec(ctx, sql, args...)
	return err
}

func (r *CommunityRepository) Update(ctx context.Context, tx pgx.Tx, c *entities.Community) error {
	ub := sqlbuilder.PostgreSQL.NewUpdateBuilder()
	ub.Update(`PUBLIC."COMMUNITY"`)
	ub.Set(
		ub.Assign("DESCRIPTION", c.Description),
		ub.Assign("RULES", c.Rules),
		ub.Assign("BANNER_URL", c.BannerURL),
		ub.Assign("PROFILE_URL", c.ProfileURL),
	)
	ub.Where(ub.Equal("ID", c.ID))

	sql, args := ub.Build()
	_, err := tx.Exec(ctx, sql, args...)
	return err
}

func (r *CommunityRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Community, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("ID", "NAME", "DESCRIPTION", "RULES", "BANNER_URL", "PROFILE_URL", "CREATED_BY", "ACTIVE", "CREATED_AT")
	sb.From(`PUBLIC."COMMUNITY"`)
	sb.Where(sb.Equal("ID", id))

	sql, args := sb.Build()
	row := r.pool.QueryRow(ctx, sql, args...)

	var c entities.Community
	err := row.Scan(&c.ID, &c.Name, &c.Description, &c.Rules, &c.BannerURL, &c.ProfileURL, &c.CreatedBy, &c.Active, &c.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *CommunityRepository) SearchPaginated(ctx context.Context, queryStr string, limit int, offset int) ([]entities.Community, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("ID", "NAME", "DESCRIPTION", "RULES", "BANNER_URL", "PROFILE_URL", "CREATED_BY", "ACTIVE", "CREATED_AT")
	sb.From(`PUBLIC."COMMUNITY"`)

	searchStr := "%" + queryStr + "%"
	sb.Where(sb.Or(
		sb.Like("NAME", searchStr),
		sb.Like("DESCRIPTION", searchStr),
	))
	sb.OrderBy("CREATED_AT").Desc()
	sb.Limit(limit).Offset(offset)

	sql, args := sb.Build()

	sql = strings.ReplaceAll(sql, "LIKE", "ILIKE")

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var communities []entities.Community
	for rows.Next() {
		var c entities.Community
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.Rules, &c.BannerURL, &c.ProfileURL, &c.CreatedBy, &c.Active, &c.CreatedAt); err != nil {
			return nil, err
		}
		communities = append(communities, c)
	}
	return communities, nil
}
