package pullrequest

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/avito/internship/pr-service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PullRequestRepository struct {
	pgxpool *pgxpool.Pool
}

func NewPullRequestRepository(pgxpool *pgxpool.Pool) *PullRequestRepository {
	return &PullRequestRepository{pgxpool: pgxpool}
}

func (r *PullRequestRepository) Get(ctx context.Context, id string) (*model.PullRequest, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sql, args, err := psql.Select(
		"pull_request_id",
		"pull_request_name",
		"author_id",
		"status",
		"created_at",
		"merged_at",
	).
		From("pull_requests").
		Where(squirrel.Eq{"pull_request_id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build get pull request query: %w", err)
	}

	var pr model.PullRequest
	err = r.pgxpool.QueryRow(ctx, sql, args...).Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrPullRequestNotFound
		}
		return nil, fmt.Errorf("failed to get pull request: %w", err)
	}

	prs := []model.PullRequest{pr}
	updatedPRs, err := r.attachReviewersToPullRequests(ctx, prs)
	if err != nil {
		return nil, err
	}

	return &updatedPRs[0], nil
}

func (r *PullRequestRepository) Create(ctx context.Context, pr *model.PullRequest) (*model.PullRequest, error) {
	tx, err := r.pgxpool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	sql, args, err := squirrel.Insert("pull_requests").
		Columns("pull_request_id", "pull_request_name", "author_id", "status", "created_at").
		Values(pr.PullRequestID, pr.PullRequestName, pr.AuthorID, pr.Status, pr.CreatedAt).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build insert pull request query: %w", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, model.ErrPullRequestExists
		}
		return nil, fmt.Errorf("failed to insert pull request: %w", err)
	}

	if len(pr.AssignedReviewers) > 0 {
		builder := squirrel.Insert("reviewers").
			Columns("pull_request_id", "user_id").
			PlaceholderFormat(squirrel.Dollar)
		for _, reviewerID := range pr.AssignedReviewers {
			builder = builder.Values(pr.PullRequestID, reviewerID)
		}
		sql, args, err := builder.ToSql()
		if err != nil {
			return nil, fmt.Errorf("failed to build insert reviewers query: %w", err)
		}
		_, err = tx.Exec(ctx, sql, args...)
		if err != nil {
			return nil, fmt.Errorf("failed to insert reviewers: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return pr, nil
}

func (r *PullRequestRepository) Merge(ctx context.Context, id string) (*model.PullRequest, error) {
	sql, args, err := squirrel.Update("pull_requests").
		Set("status", model.PullRequestMerged).
		Set("merged_at", time.Now()).
		Where(squirrel.Eq{"pull_request_id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build merge pull request query: %w", err)
	}

	result, err := r.pgxpool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to merge pull request: %w", err)
	}

	if result.RowsAffected() == 0 {
		return nil, model.ErrPullRequestNotFound
	}

	return r.Get(ctx, id)
}

func (r *PullRequestRepository) UpdateReviewer(ctx context.Context, pullRequestID, oldReviewerID, newReviewerID string) error {
	sql, args, err := squirrel.Update("reviewers").
		Set("user_id", newReviewerID).
		Where(squirrel.Eq{
			"pull_request_id": pullRequestID,
			"user_id":         oldReviewerID,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update reviewer query: %w", err)
	}

	result, err := r.pgxpool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to update reviewer: %w", err)
	}

	if result.RowsAffected() == 0 {
		return model.ErrReviewerNotAssigned
	}

	return nil
}

func (r *PullRequestRepository) attachReviewersToPullRequests(ctx context.Context, prs []model.PullRequest) ([]model.PullRequest, error) {
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

	rowsRev, err := r.pgxpool.Query(ctx, sqlRev, argsRev...)
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
