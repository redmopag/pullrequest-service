package dto

type PullRequestDTO struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	CreatedAt         string   `json:"createdAt"`
	MergedAt          string   `json:"mergedAt,omitempty"`
}

type CreatePullRequestRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type MergePullRequestRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

type ReassignPullRequestRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}
