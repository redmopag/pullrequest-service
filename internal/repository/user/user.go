package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/avito/internship/pr-service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pgxPool *pgxpool.Pool
}

func NewUserRepository(pgxPool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pgxPool: pgxPool}
}

func (r *UserRepository) UpdateStatus(ctx context.Context, userId string, isActive bool) (*model.User, error) {
	quary, args, err := squirrel.Update("users").
		Set("is_active", isActive).
		Where(squirrel.Eq{"user_id": userId}).
		Suffix("RETURNING user_id, username, team_name, is_active").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build SQL query: %w", err)
	}

	var user model.User
	err = r.pgxPool.QueryRow(ctx, quary, args...).Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("update user status: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUsersPullRequests(ctx context.Context, userId string) ([]model.PullRequest, error) {
	sqlPr, argsPr, err := squirrel.Select(
		"pull_request_id",
		"pull_request_name",
		"author_id",
		"status",
		"created_at",
		"merged_at",
	).
		From("pull_requests").
		Where(squirrel.Eq{"author_id": userId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build pull requests query: %w", err)
	}

	rows, err := r.pgxPool.Query(ctx, sqlPr, argsPr...)
	if err != nil {
		return nil, fmt.Errorf("query pull requests: %w", err)
	}
	defer rows.Close()

	var pullRequests []model.PullRequest
	for rows.Next() {
		var pr model.PullRequest
		if err := rows.Scan(
			&pr.PullRequestID,
			&pr.PullRequestName,
			&pr.AuthorID,
			&pr.Status,
			&pr.CreatedAt,
			&pr.MergedAt,
		); err != nil {
			return nil, fmt.Errorf("scan pull request: %w", err)
		}
		pullRequests = append(pullRequests, pr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	if len(pullRequests) == 0 {
		return []model.PullRequest{}, nil
	}

	updatedPRs, err := r.attachReviewersToPullRequests(ctx, pullRequests)
	if err != nil {
		return nil, err
	}

	return updatedPRs, nil
}

func (r *UserRepository) attachReviewersToPullRequests(ctx context.Context, prs []model.PullRequest) ([]model.PullRequest, error) {
	if len(prs) == 0 {
		return prs, nil
	}

	prIDs := make([]string, 0, len(prs))
	for _, pr := range prs {
		prIDs = append(prIDs, pr.PullRequestID)
	}

	sqlRev, argsRev, err := squirrel.Select("pull_request_id", "user_id").
		From("reviewers").
		Where(squirrel.Eq{"pull_request_id": prIDs}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build reviewers query: %w", err)
	}

	rowsRev, err := r.pgxPool.Query(ctx, sqlRev, argsRev...)
	if err != nil {
		return nil, fmt.Errorf("query reviewers: %w", err)
	}
	defer rowsRev.Close()

	reviewersMap := make(map[string][]string)
	for rowsRev.Next() {
		var prID, userID string
		if err := rowsRev.Scan(&prID, &userID); err != nil {
			return nil, fmt.Errorf("scan reviewers: %w", err)
		}
		reviewersMap[prID] = append(reviewersMap[prID], userID)
	}

	if err := rowsRev.Err(); err != nil {
		return nil, fmt.Errorf("reviewers rows error: %w", err)
	}

	for i := range prs {
		if revs, ok := reviewersMap[prs[i].PullRequestID]; ok {
			prs[i].AssignedReviewers = revs
		} else {
			prs[i].AssignedReviewers = []string{}
		}
	}

	return prs, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userId string) (*model.User, error) {
	query, args, err := squirrel.Select("user_id", "username", "team_name", "is_active").
		From("users").
		Where(squirrel.Eq{"user_id": userId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build SQL query: %w", err)
	}

	var user model.User
	err = r.pgxPool.QueryRow(ctx, query, args...).Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	return &user, nil
}
