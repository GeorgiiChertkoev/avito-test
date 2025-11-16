package postgres

import (
	"context"
	"database/sql"
	"errors"
	"pr-reviewer/internal/domain/user"
)

type UserRepo struct {
	tx *sql.Tx
}

func NewUserRepo(tx *sql.Tx) *UserRepo {
	return &UserRepo{tx: tx}
}

func (r *UserRepo) UpsertUser(ctx context.Context, usr *user.User) error {
	query := `
		INSERT INTO "user" (user_id, username, team_name, is_active)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id)
		DO UPDATE SET 
			username = EXCLUDED.username,
			team_name = EXCLUDED.team_name,
			is_active = EXCLUDED.is_active
	`

	_, err := r.tx.ExecContext(ctx, query,
		usr.UserID,
		usr.Username,
		usr.TeamName,
		usr.IsActive,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return user.ErrNotFound
	}

	return err
}

func (r *UserRepo) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	query := `
		SELECT user_id, username, team_name, is_active
		FROM "user"
		WHERE user_id = $1
	`

	row := r.tx.QueryRowContext(ctx, query, id)

	var u user.User
	err := row.Scan(
		&u.UserID,
		&u.Username,
		&u.TeamName,
		&u.IsActive,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, user.ErrNotFound
	}

	return &u, err
}

func (r *UserRepo) SetIsActive(ctx context.Context, id string, active bool) (*user.User, error) {
	query := `
		UPDATE "user"
		SET is_active = $2
		WHERE user_id = $1
		RETURNING user_id, username, team_name, is_active
	`

	var u user.User
	err := r.tx.QueryRowContext(ctx, query, id, active).Scan(
		&u.UserID,
		&u.Username,
		&u.TeamName,
		&u.IsActive,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // user not found
	}

	return &u, err
}
