package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, id uuid.UUID) error {
	query := `INSERT INTO "USER" (ID) VALUES ($1)`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to insert user %s: %w", id, err)
	}
	return nil
}
