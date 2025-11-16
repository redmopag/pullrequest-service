package user

import (
	"context"

	"github.com/avito/internship/pr-service/internal/model"
)

type UserRepository interface {
	UpdateStatus(ctx context.Context, userId string, isActive bool) (*model.User, error)
	GetUserByID(ctx context.Context, userId string) (*model.User, error)
}
