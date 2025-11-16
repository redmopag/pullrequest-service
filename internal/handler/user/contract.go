package user

import (
	"context"

	"github.com/avito/internship/pr-service/internal/model"
)

type UserService interface {
	SetActiveStatus(ctx context.Context, userId string, isActive bool) (*model.User, error)
	GetUsersPullRequests(ctx context.Context, userId string) ([]model.PullRequest, error)
	GetUserByID(ctx context.Context, userId string) (*model.User, error)
}
