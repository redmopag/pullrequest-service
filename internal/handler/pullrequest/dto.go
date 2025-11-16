package pullrequest

import (
	"time"

	"github.com/avito/internship/pr-service/internal/model"
)

type PullRequestDTO struct {
	PullRequestID     string                  `json:"pull_request_id"`
	PullRequestName   string                  `json:"pull_request_name"`
	AuthorID          string                  `json:"author_id"`
	Status            model.PullRequestStatus `json:"status"`
	AssignedReviewers []string                `json:"assigned_reviewers"`
	MergedAt          *time.Time              `json:"merged_at,omitempty"`
}

type ReassignPullRequestRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_reviewer_id"`
}

type ReassignResponseDTO struct {
	PullRequest PullRequestDTO `json:"pr"`
	ReplacedBy  string         `json:"replaced_by"`
}

type CreatePullRequestRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type MergePullRequestRequest struct {
	PullRequestID string `json:"pull_request_id"`
}
