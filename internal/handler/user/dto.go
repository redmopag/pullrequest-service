package user

import "github.com/avito/internship/pr-service/internal/model"

type SetIsActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type UserResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type PullRequestShort struct {
	PullRequestID   string                  `json:"pull_request_id"`
	PullRequestName string                  `json:"pull_request_name"`
	AuthorID        string                  `json:"author_id"`
	Status          model.PullRequestStatus `json:"status"`
}
