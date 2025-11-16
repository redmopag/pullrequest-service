package team

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/avito/internship/pr-service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepository struct {
	pgxpool *pgxpool.Pool
}

func NewTeamRepository(pgxpool *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{pgxpool: pgxpool}
}

func (r *TeamRepository) CreateTeam(ctx context.Context, team *model.Team) (*model.Team, error) {
	tx, err := r.pgxpool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	err = insertTeam(ctx, team, tx)
	if err != nil {
		return nil, err
	}

	err = insertUsers(ctx, team, tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return team, nil
}

func insertTeam(ctx context.Context, team *model.Team, tx pgx.Tx) error {
	sql, args, err := squirrel.Insert("teams").
		Columns("team_name").
		Values(team.TeamName).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.ErrTeamExists
		}
	}
	return nil
}

func insertUsers(ctx context.Context, team *model.Team, tx pgx.Tx) error {
	if len(team.Users) > 0 {
		builder := squirrel.Insert("users").
			Columns("team_name", "user_id", "username", "is_active").
			PlaceholderFormat(squirrel.Dollar)
		for _, user := range team.Users {
			builder = builder.Values(team.TeamName, user.UserID, user.Username, user.IsActive)
		}
		builder = builder.Suffix("ON CONFLICT (user_id) DO UPDATE SET username = EXCLUDED.username, team_name = EXCLUDED.team_name, is_active = EXCLUDED.is_active")
		sql, args, err := builder.ToSql()
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, sql, args...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *TeamRepository) GetTeamByName(ctx context.Context, teamName string) (*model.Team, error) {
	sqlBuilder, args, err := squirrel.Select("t.team_name", "u.user_id", "u.username", "u.is_active").
		From("teams t").
		LeftJoin("users u ON t.team_name = u.team_name").
		Where(squirrel.Eq{"t.team_name": teamName}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.pgxpool.Query(ctx, sqlBuilder, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var team model.Team
	teamFound := false
	for rows.Next() {
		var userID sql.NullString
		var username sql.NullString
		var isActive sql.NullBool

		var scannedTeamName string

		err := rows.Scan(&scannedTeamName, &userID, &username, &isActive)
		if err != nil {
			return nil, err
		}

		if !teamFound {
			team.TeamName = scannedTeamName
			teamFound = true
		}

		if userID.Valid {
			user := model.User{
				UserID:   userID.String,
				Username: username.String,
				IsActive: isActive.Bool,
				TeamName: team.TeamName,
			}
			team.Users = append(team.Users, user)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if !teamFound {
		return nil, model.ErrTeamNotFound
	}
	return &team, nil
}
