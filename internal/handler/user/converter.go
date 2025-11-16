package user

import "github.com/avito/internship/pr-service/internal/model"

func ToUserResponse(user *model.User) *UserResponse {
	return &UserResponse{
		UserID:   user.UserID,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
}

func ToPullRequestShort(pr *model.PullRequest) *PullRequestShort {
	return &PullRequestShort{
		PullRequestID:   pr.PullRequestID,
		PullRequestName: pr.PullRequestName,
		AuthorID:        pr.AuthorID,
		Status:          pr.Status,
	}
}
