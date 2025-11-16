package pullrequest

import (
	"context"

	"github.com/avito/internship/pr-service/internal/model"
)

type PullRequestRepository interface {
	Get(ctx context.Context, id string) (*model.PullRequest, error)
	Create(ctx context.Context, pr *model.PullRequest) (*model.PullRequest, error)
	Merge(ctx context.Context, id string) (*model.PullRequest, error)
	UpdateReviewer(ctx context.Context, pullRequestID, oldReviewerID, newReviewerID string) error
}

type UserService interface {
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
}

type TeamService interface {
	GetTeamByName(ctx context.Context, teamName string) (*model.Team, error)
}
