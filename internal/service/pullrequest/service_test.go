package pullrequest_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/avito/internship/pr-service/internal/model"
	"github.com/avito/internship/pr-service/internal/service/pullrequest"
	"github.com/avito/internship/pr-service/internal/service/pullrequest/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPullRequestService_Create(t *testing.T) {
	ctx := context.Background()
	author := &model.User{UserID: "author1", TeamName: "team-a", IsActive: true}
	team := &model.Team{
		TeamName: "team-a",
		Users: []model.User{
			*author,
			{UserID: "user2", IsActive: true},
			{UserID: "user3", IsActive: true},
			{UserID: "user4", IsActive: false},
		},
	}

	tests := []struct {
		name          string
		setupMocks    func(*mocks.PullRequestRepository, *mocks.UserService, *mocks.TeamService)
		prID          string
		prName        string
		authorID      string
		expectedError string
	}{
		{
			name:     "Success",
			prID:     "pr1",
			prName:   "New Feature",
			authorID: "author1",
			setupMocks: func(pullRequestRepository *mocks.PullRequestRepository, userService *mocks.UserService, teamService *mocks.TeamService) {
				userService.On("GetUserByID", ctx, "author1").Return(author, nil).Once()
				teamService.On("GetTeamByName", ctx, "team-a").Return(team, nil).Once()
				pullRequestRepository.On("Create", ctx, mock.AnythingOfType("*model.PullRequest")).Run(func(args mock.Arguments) {
					pr := args.Get(1).(*model.PullRequest)
					assert.Equal(t, "pr1", pr.PullRequestID)
					assert.Equal(t, "author1", pr.AuthorID)
					assert.Equal(t, model.PullRequestOpen, pr.Status)
					assert.Len(t, pr.AssignedReviewers, 2)
					assert.NotContains(t, pr.AssignedReviewers, "author1")
					assert.NotContains(t, pr.AssignedReviewers, "user4")
				}).Return(&model.PullRequest{}, nil).Once()
			},
			expectedError: "",
		},
		{
			name:     "Error user not found",
			prID:     "pr1",
			prName:   "New Feature",
			authorID: "author1",
			setupMocks: func(pullRequestRepository *mocks.PullRequestRepository, userService *mocks.UserService, teamService *mocks.TeamService) {
				userService.On("GetUserByID", ctx, "author1").Return(nil, model.ErrUserNotFound).Once()
			},
			expectedError: "get author: user not found",
		},
		{
			name:     "Error team not found",
			prID:     "pr1",
			prName:   "New Feature",
			authorID: "author1",
			setupMocks: func(pullRequestRepository *mocks.PullRequestRepository, userService *mocks.UserService, teamService *mocks.TeamService) {
				userService.On("GetUserByID", ctx, "author1").Return(author, nil).Once()
				teamService.On("GetTeamByName", ctx, "team-a").Return(nil, model.ErrTeamNotFound).Once()
			},
			expectedError: "get team: team not found",
		},
		{
			name:     "Error on create PR",
			prID:     "pr1",
			prName:   "New Feature",
			authorID: "author1",
			setupMocks: func(pullRequestRepository *mocks.PullRequestRepository, userService *mocks.UserService, teamService *mocks.TeamService) {
				userService.On("GetUserByID", ctx, "author1").Return(author, nil).Once()
				teamService.On("GetTeamByName", ctx, "team-a").Return(team, nil).Once()
				pullRequestRepository.On("Create", ctx, mock.AnythingOfType("*model.PullRequest")).Return(nil, errors.New("db error")).Once()
			},
			expectedError: "create PR: db error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pullRequestRepository := new(mocks.PullRequestRepository)
			userService := new(mocks.UserService)
			teamService := new(mocks.TeamService)
			test.setupMocks(pullRequestRepository, userService, teamService)

			s := pullrequest.NewPullRequestService(pullRequestRepository, userService, teamService)
			_, err := s.Create(ctx, test.prID, test.prName, test.authorID)

			if test.expectedError != "" {
				assert.EqualError(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}

			pullRequestRepository.AssertExpectations(t)
			userService.AssertExpectations(t)
			teamService.AssertExpectations(t)
		})
	}
}

func TestPullRequestService_Merge(t *testing.T) {
	ctx := context.Background()
	prOpen := &model.PullRequest{PullRequestID: "pr1", Status: model.PullRequestOpen}
	prMerged := &model.PullRequest{PullRequestID: "pr1", Status: model.PullRequestMerged}

	tests := []struct {
		name          string
		setupMocks    func(*mocks.PullRequestRepository)
		prID          string
		expectedError string
	}{
		{
			name: "Success - Merge open PR",
			prID: "pr1",
			setupMocks: func(pullRequestRepository *mocks.PullRequestRepository) {
				pullRequestRepository.On("Get", ctx, "pr1").Return(prOpen, nil).Once()
				pullRequestRepository.On("Merge", ctx, "pr1").Return(prMerged, nil).Once()
			},
			expectedError: "",
		},
		{
			name: "Success - PR already merged",
			prID: "pr1",
			setupMocks: func(pullRequestRepository *mocks.PullRequestRepository) {
				pullRequestRepository.On("Get", ctx, "pr1").Return(prMerged, nil).Once()
			},
			expectedError: "",
		},
		{
			name: "Error - PR not found on Get",
			prID: "pr1",
			setupMocks: func(pullRequestRepository *mocks.PullRequestRepository) {
				pullRequestRepository.On("Get", ctx, "pr1").Return(nil, model.ErrPullRequestNotFound).Once()
			},
			expectedError: "get PR to merge: pull request not found",
		},
		{
			name: "Error - PR not found on Merge",
			prID: "pr1",
			setupMocks: func(pullRequestRepository *mocks.PullRequestRepository) {
				pullRequestRepository.On("Get", ctx, "pr1").Return(prOpen, nil).Once()
				pullRequestRepository.On("Merge", ctx, "pr1").Return(nil, model.ErrPullRequestNotFound).Once()
			},
			expectedError: "merge PR: pull request not found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pullRequestRepository := new(mocks.PullRequestRepository)
			test.setupMocks(pullRequestRepository)

			s := pullrequest.NewPullRequestService(pullRequestRepository, nil, nil)
			_, err := s.Merge(ctx, test.prID)

			if test.expectedError != "" {
				assert.EqualError(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}
			pullRequestRepository.AssertExpectations(t)
		})
	}
}

func TestPullRequestService_Reassign(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	prOpen := &model.PullRequest{
		PullRequestID:     "pr1",
		Status:            model.PullRequestOpen,
		AuthorID:          "author1",
		AssignedReviewers: []string{"old_reviewer", "other_reviewer"},
	}
	prMerged := &model.PullRequest{PullRequestID: "pr1", Status: model.PullRequestMerged, MergedAt: &now}
	oldReviewer := &model.User{UserID: "old_reviewer", TeamName: "team-a", IsActive: true}
	team := &model.Team{
		TeamName: "team-a",
		Users: []model.User{
			{UserID: "author1"},
			*oldReviewer,
			{UserID: "other_reviewer"},
			{UserID: "new_reviewer", IsActive: true},
			{UserID: "inactive_reviewer", IsActive: false},
		},
	}
	prWithNewReviewer := &model.PullRequest{
		PullRequestID:     "pr1",
		Status:            model.PullRequestOpen,
		AuthorID:          "author1",
		AssignedReviewers: []string{"new_reviewer", "other_reviewer"},
	}

	tests := []struct {
		name          string
		setupMocks    func(*mocks.PullRequestRepository, *mocks.UserService, *mocks.TeamService)
		prID          string
		oldReviewerID string
		expectedError string
	}{
		{
			name:          "Success",
			prID:          "pr1",
			oldReviewerID: "old_reviewer",
			setupMocks: func(pullRequestRepository *mocks.PullRequestRepository, userService *mocks.UserService, teamService *mocks.TeamService) {
				pullRequestRepository.On("Get", ctx, "pr1").Return(prOpen, nil).Once()
				userService.On("GetUserByID", ctx, "old_reviewer").Return(oldReviewer, nil).Once()
				teamService.On("GetTeamByName", ctx, "team-a").Return(team, nil).Once()
				pullRequestRepository.On("UpdateReviewer", ctx, "pr1", "old_reviewer", "new_reviewer").Return(nil).Once()
				pullRequestRepository.On("Get", ctx, "pr1").Return(prWithNewReviewer, nil).Once()
			},
			expectedError: "",
		},
		{
			name:          "Error - PR is merged",
			prID:          "pr1",
			oldReviewerID: "old_reviewer",
			setupMocks: func(pullRequestRepository *mocks.PullRequestRepository, userService *mocks.UserService, teamService *mocks.TeamService) {
				pullRequestRepository.On("Get", ctx, "pr1").Return(prMerged, nil).Once()
			},
			expectedError: "cannot reassign on merged PR",
		},
		{
			name:          "Error - Reviewer not assigned",
			prID:          "pr1",
			oldReviewerID: "not_assigned_reviewer",
			setupMocks: func(pullRequestRepository *mocks.PullRequestRepository, userService *mocks.UserService, teamService *mocks.TeamService) {
				pullRequestRepository.On("Get", ctx, "pr1").Return(prOpen, nil).Once()
			},
			expectedError: "reviewer is not assigned to this pull request",
		},
		{
			name:          "Error - No candidates for reassign",
			prID:          "pr1",
			oldReviewerID: "old_reviewer",
			setupMocks: func(pullRequestRepository *mocks.PullRequestRepository, userService *mocks.UserService, teamService *mocks.TeamService) {
				pullRequestRepository.On("Get", ctx, "pr1").Return(prOpen, nil).Once()
				userService.On("GetUserByID", ctx, "old_reviewer").Return(oldReviewer, nil).Once()
				teamWithoutCandidates := &model.Team{
					TeamName: "team-a",
					Users: []model.User{
						{UserID: "author1"},
						*oldReviewer,
						{UserID: "other_reviewer"},
					},
				}
				teamService.On("GetTeamByName", ctx, "team-a").Return(teamWithoutCandidates, nil).Once()
			},
			expectedError: "no active replacement candidate in team",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pullRequestRepository := new(mocks.PullRequestRepository)
			userService := new(mocks.UserService)
			teamService := new(mocks.TeamService)
			test.setupMocks(pullRequestRepository, userService, teamService)

			s := pullrequest.NewPullRequestService(pullRequestRepository, userService, teamService)
			_, err := s.Reassign(ctx, test.prID, test.oldReviewerID)

			if test.expectedError != "" {
				assert.EqualError(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}

			pullRequestRepository.AssertExpectations(t)
			userService.AssertExpectations(t)
			teamService.AssertExpectations(t)
		})
	}
}

func TestPullRequestService_GetPullRequestsForUserReview(t *testing.T) {
	ctx := context.Background()
	userID := "user123"

	tests := []struct {
		name          string
		setupMocks    func(*mocks.PullRequestRepository, *mocks.UserService)
		userID        string
		expectedPRs   []model.PullRequest
		expectedError error
	}{
		{
			name:   "Error - User not found",
			userID: userID,
			setupMocks: func(pullRequestRepository *mocks.PullRequestRepository, userService *mocks.UserService) {
				userService.On("GetUserByID", ctx, userID).Return(nil, model.ErrUserNotFound).Once()
				pullRequestRepository.AssertNotCalled(t, "GetUsersPullRequests", ctx, userID)
			},
			expectedPRs:   nil,
			expectedError: model.ErrUserNotFound,
		},
		{
			name:   "Error - UserService returns generic error",
			userID: userID,
			setupMocks: func(pullRequestRepository *mocks.PullRequestRepository, userService *mocks.UserService) {
				userService.On("GetUserByID", ctx, userID).Return(nil, errors.New("some user service error")).Once()
				pullRequestRepository.AssertNotCalled(t, "GetUsersPullRequests", ctx, userID)
			},
			expectedPRs:   nil,
			expectedError: errors.New("some user service error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pullRequestRepository := new(mocks.PullRequestRepository)
			userService := new(mocks.UserService)
			s := pullrequest.NewPullRequestService(pullRequestRepository, userService, nil)

			test.setupMocks(pullRequestRepository, userService)

			actualPRs, actualError := s.GetPullRequestsForUserReview(ctx, test.userID)

			assert.Equal(t, test.expectedPRs, actualPRs)
			assert.Equal(t, test.expectedError, actualError)

			pullRequestRepository.AssertExpectations(t)
			userService.AssertExpectations(t)
		})
	}
}
