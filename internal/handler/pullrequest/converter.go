package pullrequest

import "github.com/avito/internship/pr-service/internal/model"

func ToPullRequestDTO(pr *model.PullRequest) *PullRequestDTO {
	return &PullRequestDTO{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            pr.Status,
		AssignedReviewers: pr.AssignedReviewers,
		MergedAt:          pr.MergedAt,
	}
}
