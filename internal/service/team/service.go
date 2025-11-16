package team

import (
	"context"

	"github.com/avito/internship/pr-service/internal/model"
)

type TeamService struct {
	taemRepository TeamRepository
}

func NewTeamService(repository TeamRepository) *TeamService {
	return &TeamService{
		taemRepository: repository,
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, team *model.Team) (*model.Team, error) {
	return s.taemRepository.CreateTeam(ctx, team)
}

func (s *TeamService) GetTeamByName(ctx context.Context, teamName string) (*model.Team, error) {
	return s.taemRepository.GetTeamByName(ctx, teamName)
}
