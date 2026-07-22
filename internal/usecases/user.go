package usecases

import (
	"context"

	"github.com/google/uuid"
)

type UserUseCase struct {
	userRepo UserRepository
}

func NewUserUseCase(userRepo UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, id uuid.UUID) error {
	return uc.userRepo.Create(ctx, id)
}
