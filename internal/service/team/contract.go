package team

import (
	"context"

	"github.com/avito/internship/pr-service/internal/model"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *model.Team) (*model.Team, error)
	GetTeamByName(ctx context.Context, teamName string) (*model.Team, error)
}
