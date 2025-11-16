package model

type DomainError struct {
	Code    string
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

var (
	ErrTeamExists           = &DomainError{Code: "TEAM_EXISTS", Message: "team_name already exists"}
	ErrTeamNotFound         = &DomainError{Code: "NOT_FOUND", Message: "team not found"}
	ErrUserNotFound         = &DomainError{Code: "NOT_FOUND", Message: "user not found"}
	ErrPullRequestNotFound  = &DomainError{Code: "NOT_FOUND", Message: "pull request not found"}
	ErrPullRequestExists    = &DomainError{Code: "PR_EXISTS", Message: "pull request already exists"}
	ErrPRMerged             = &DomainError{Code: "PR_MERGED", Message: "cannot reassign on merged PR"}
	ErrReviewerNotAssigned  = &DomainError{Code: "NOT_ASSIGNED", Message: "reviewer is not assigned to this pull request"}
	ErrNoReviewerCandidates = &DomainError{Code: "NO_CANDIDATE", Message: "no active replacement candidate in team"}
	ErrBadJSONRequest       = &DomainError{Code: "BAD_REQUEST", Message: "bad json request"}
	ErrTeamNoUsers          = &DomainError{Code: "BAD_REQUEST", Message: "team must have at least one user"}
)
