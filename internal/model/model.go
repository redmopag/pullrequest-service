package model

import (
	"time"
)

type User struct {
	UserID   string
	Username string
	TeamName string
	IsActive bool
}

type Team struct {
	TeamName string
	Users    []User
}

type PullRequest struct {
	PullRequestID     string
	PullRequestName   string
	AuthorID          string
	Status            PullRequestStatus
	AssignedReviewers []string
	CreatedAt         time.Time
	MergedAt          time.Time
}

type PullRequestStatus string

const (
	PullRequestMerged PullRequestStatus = "MERGED"
	PullRequestOpen   PullRequestStatus = "OPEN"
)

type ReassignResponse struct {
	PullRequest *PullRequest
	ReplacedBy  string
}
