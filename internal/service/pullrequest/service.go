package pullrequest

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/avito/internship/pr-service/internal/model"
)

type PullRequestService struct {
	pullRequestRepository PullRequestRepository
	userService           UserService
	teamService           TeamService
}

func NewPullRequestService(pullRequestRepository PullRequestRepository, userService UserService, teamService TeamService) *PullRequestService {
	return &PullRequestService{
		pullRequestRepository: pullRequestRepository,
		userService:           userService,
		teamService:           teamService,
	}
}

func (s *PullRequestService) Create(ctx context.Context, prID, prName, authorID string) (*model.PullRequest, error) {
	author, err := s.userService.GetUserByID(ctx, authorID)
	if err != nil {
		return nil, fmt.Errorf("get author: %w", err)
	}

	team, err := s.teamService.GetTeamByName(ctx, author.TeamName)
	if err != nil {
		return nil, fmt.Errorf("get team: %w", err)
	}

	reviewers := s.selectReviewers(authorID, team.Users)

	pr := &model.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   prName,
		AuthorID:          authorID,
		Status:            model.PullRequestOpen,
		AssignedReviewers: reviewers,
		CreatedAt:         time.Now(),
	}

	createdPR, err := s.pullRequestRepository.Create(ctx, pr)
	if err != nil {
		return nil, fmt.Errorf("create PR: %w", err)
	}

	return createdPR, nil
}

func (s *PullRequestService) Merge(ctx context.Context, prID string) (*model.PullRequest, error) {
	pr, err := s.pullRequestRepository.Get(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("get PR to merge: %w", err)
	}

	if pr.Status == model.PullRequestMerged {
		return pr, nil
	}

	mergedPR, err := s.pullRequestRepository.Merge(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("merge PR: %w", err)
	}

	return mergedPR, nil
}

func (s *PullRequestService) Reassign(ctx context.Context, prID, oldReviewerID string) (*model.ReassignResponse, error) {
	pr, err := s.pullRequestRepository.Get(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("get PR to reassign: %w", err)
	}

	if pr.Status == model.PullRequestMerged {
		return nil, model.ErrPRMerged
	}

	isAssigned := false
	for _, r := range pr.AssignedReviewers {
		if r == oldReviewerID {
			isAssigned = true
			break
		}
	}
	if !isAssigned {
		return nil, model.ErrReviewerNotAssigned
	}

	oldReviewer, err := s.userService.GetUserByID(ctx, oldReviewerID)
	if err != nil {
		return nil, fmt.Errorf("get old reviewer: %w", err)
	}

	team, err := s.teamService.GetTeamByName(ctx, oldReviewer.TeamName)
	if err != nil {
		return nil, fmt.Errorf("get team: %w", err)
	}

	newReviewerID, err := s.selectNewReviewer(pr.AuthorID, pr.AssignedReviewers, team.Users)
	if err != nil {
		return nil, err
	}

	err = s.pullRequestRepository.UpdateReviewer(ctx, prID, oldReviewerID, newReviewerID)
	if err != nil {
		return nil, fmt.Errorf("update reviewer: %w", err)
	}

	updatedPR, err := s.pullRequestRepository.Get(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("get updated PR: %w", err)
	}

	return &model.ReassignResponse{
		PullRequest: updatedPR,
		ReplacedBy:  newReviewerID,
	}, nil
}

func (s *PullRequestService) GetPullRequestsForUserReview(ctx context.Context, userId string) ([]model.PullRequest, error) {
	if _, err := s.userService.GetUserByID(ctx, userId); err != nil {
		return nil, err
	}
	return s.pullRequestRepository.GetUsersPullRequests(ctx, userId)
}

func (s *PullRequestService) selectReviewers(authorID string, users []model.User) []string {
	candidates := make([]string, 0)
	for _, user := range users {
		if user.IsActive && user.UserID != authorID {
			candidates = append(candidates, user.UserID)
		}
	}

	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	if len(candidates) > 2 {
		return candidates[:2]
	}
	return candidates
}

func (s *PullRequestService) selectNewReviewer(authorID string, currentReviewers []string, users []model.User) (string, error) {
	candidates := make([]string, 0)
	isReviewer := func(userID string) bool {
		for _, r := range currentReviewers {
			if r == userID {
				return true
			}
		}
		return false
	}

	for _, user := range users {
		if user.IsActive && user.UserID != authorID && !isReviewer(user.UserID) {
			candidates = append(candidates, user.UserID)
		}
	}

	if len(candidates) == 0 {
		return "", model.ErrNoReviewerCandidates
	}

	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	return candidates[0], nil
}
