package mocks

import (
	"context"

	"github.com/avito/internship/pr-service/internal/model"
	"github.com/stretchr/testify/mock"
)

type PullRequestRepository struct {
	mock.Mock
}

func (m *PullRequestRepository) Get(ctx context.Context, id string) (*model.PullRequest, error) {
	args := m.Called(ctx, id)
	var pullRequest *model.PullRequest
	if args.Get(0) != nil {
		pullRequest = args.Get(0).(*model.PullRequest)
	}
	return pullRequest, args.Error(1)
}

func (m *PullRequestRepository) Create(ctx context.Context, pr *model.PullRequest) (*model.PullRequest, error) {
	args := m.Called(ctx, pr)
	var createdPullRequest *model.PullRequest
	if args.Get(0) != nil {
		createdPullRequest = args.Get(0).(*model.PullRequest)
	}
	return createdPullRequest, args.Error(1)
}

func (m *PullRequestRepository) Merge(ctx context.Context, id string) (*model.PullRequest, error) {
	args := m.Called(ctx, id)
	var pullRequest *model.PullRequest
	if args.Get(0) != nil {
		pullRequest = args.Get(0).(*model.PullRequest)
	}
	return pullRequest, args.Error(1)
}

func (m *PullRequestRepository) UpdateReviewer(ctx context.Context, pullRequestID, oldReviewerID, newReviewerID string) error {
	args := m.Called(ctx, pullRequestID, oldReviewerID, newReviewerID)
	return args.Error(0)
}

func (m *PullRequestRepository) GetUsersPullRequests(ctx context.Context, userId string) ([]model.PullRequest, error) {
	args := m.Called(ctx, userId)
	var pullRequests []model.PullRequest
	if args.Get(0) != nil {
		pullRequests = args.Get(0).([]model.PullRequest)
	}
	return pullRequests, args.Error(1)
}
