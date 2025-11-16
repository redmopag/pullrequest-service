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
