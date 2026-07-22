package database

import (
	"context"

	"github.com/RedditUclaista/community-service/internal/entities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxRepository struct {
	pool *pgxpool.Pool
}

func NewOutboxRepository(pool *pgxpool.Pool) *OutboxRepository {
	return &OutboxRepository{pool: pool}
}

func (r *OutboxRepository) Insert(ctx context.Context, tx pgx.Tx, event *entities.OutboxEvent) error {
	query := `
		INSERT INTO PUBLIC."OUTBOX_EVENT" (ID, AGGREGATE_TYPE, AGGREGATE_ID, TYPE, PAYLOAD, STATUS, CREATED_AT, UPDATED_AT)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := tx.Exec(ctx, query, event.ID, event.AggregateType, event.AggregateID, event.Type, event.Payload, event.Status, event.CreatedAt, event.UpdatedAt)
	return err
}
