package user

import (
	"context"

	"github.com/avito/internship/pr-service/internal/model"
)

type UserService struct {
	userRepository UserRepository
}

func NewUserService(repostiroy UserRepository) *UserService {
	return &UserService{userRepository: repostiroy}
}

func (service *UserService) SetActiveStatus(ctx context.Context, userId string, isActive bool) (*model.User, error) {
	return service.userRepository.UpdateStatus(ctx, userId, isActive)
}

func (service *UserService) GetUsersPullRequests(ctx context.Context, userId string) ([]model.PullRequest, error) {
	return service.userRepository.GetUsersPullRequests(ctx, userId)
}

func (service *UserService) GetUserByID(ctx context.Context, userId string) (*model.User, error) {
	return service.userRepository.GetUserByID(ctx, userId)
}
