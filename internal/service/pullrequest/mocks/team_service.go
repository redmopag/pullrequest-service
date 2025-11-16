package mocks

import (
	"context"

	"github.com/avito/internship/pr-service/internal/model"
	"github.com/stretchr/testify/mock"
)

type TeamService struct {
	mock.Mock
}

func (m *TeamService) GetTeamByName(ctx context.Context, teamName string) (*model.Team, error) {
	args := m.Called(ctx, teamName)
	var team *model.Team
	if args.Get(0) != nil {
		team = args.Get(0).(*model.Team)
	}
	return team, args.Error(1)
}
