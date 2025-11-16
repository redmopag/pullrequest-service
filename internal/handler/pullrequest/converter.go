package pullrequest

import "github.com/avito/internship/pr-service/internal/model"

func ToPullRequestDTO(pr *model.PullRequest) *PullRequestDTO {
	return &PullRequestDTO{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            string(pr.Status),
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}
