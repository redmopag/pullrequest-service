package mocks

import (
	"context"

	"github.com/avito/internship/pr-service/internal/model"
	"github.com/stretchr/testify/mock"
)

type UserService struct {
	mock.Mock
}

func (m *UserService) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	args := m.Called(ctx, userID)
	var user *model.User
	if args.Get(0) != nil {
		user = args.Get(0).(*model.User)
	}
	return user, args.Error(1)
}
