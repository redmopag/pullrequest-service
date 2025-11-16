package pullrequest

import (
	"context"

	"github.com/avito/internship/pr-service/internal/model"
)

type PullRequestService interface {
	Create(ctx context.Context, prID, prName, authoЯrID string) (*model.PullRequest, error)
	Merge(ctx context.Context, prID string) (*model.PullRequest, error)
	Reassign(ctx context.Context, prID, oldReviewerID string) (*model.ReassignResponse, error)
}
