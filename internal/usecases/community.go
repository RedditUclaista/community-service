package usecases

import (
	"database/sql"
		"github.com/google/uuid"
	"github.com/RedditUclaista/community-service/internal/entities"
)

type CommunityUseCase struct {
	db *sql.DB
}

func NewCommunityUseCase(db *sql.DB) *CommunityUseCase {
	return &CommunityUseCase{db: db}
}

func (uc *CommunityUseCase) Create(name, description, createdBy string) (*entities.Community, error) {
	id := uuid.New()
	query := `INSERT INTO communities (id, name, description, created_by, active, created_at) 
	          VALUES ($1, $2, $3, $4, true, NOW()) RETURNING created_at`
	
	comm := &entities.Community{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedBy:   createdBy,
		Active:      true,
	}
	
	err := uc.db.QueryRow(query, id, name, description, createdBy).Scan(&comm.CreatedAt)
	if err != nil {
		return nil, err
	}
	
	return comm, nil
}

func (uc *CommunityUseCase) ListActive() ([]entities.Community, error) {
	query := `SELECT id, name, description, created_by, active, created_at FROM communities WHERE active = true`
	rows, err := uc.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var list []entities.Community
	for rows.Next() {
		var c entities.Community
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedBy, &c.Active, &c.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}
