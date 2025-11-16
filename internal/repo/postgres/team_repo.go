package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"pr-reviewer/internal/domain/team"
)

type TeamRepo struct {
	tx *sql.Tx
}

func NewTeamRepo(tx *sql.Tx) *TeamRepo {
	return &TeamRepo{tx: tx}
}

func (r *TeamRepo) CreateTeam(ctx context.Context, t *team.Team) error {
	//check that doesn't exist yet
	_, err := r.GetTeamByName(ctx, t.TeamName)
	if !errors.Is(err, team.ErrNotFound) {
		return team.ErrTeamExists
	}

	queryTeam := `
       INSERT INTO team (team_name)
       VALUES ($1)
   `
	if _, err := r.tx.ExecContext(ctx, queryTeam, t.TeamName); err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	for _, m := range t.Members {
		queryUser := `
			INSERT INTO "user" (user_id, username, team_name, is_active)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (user_id)
			DO UPDATE SET 
				username = EXCLUDED.username,
				team_name = EXCLUDED.team_name,
				is_active = EXCLUDED.is_active
       `
		_, err := r.tx.ExecContext(ctx, queryUser,
			m.UserID,
			m.Username,
			t.TeamName,
			m.IsActive,
		)
		if err != nil {
			return fmt.Errorf("failed to insert user %q: %w", m.Username, err)
		}
	}

	return nil
}

func (r *TeamRepo) GetTeamByName(ctx context.Context, name string) (*team.Team, error) {
	var teamID string

	queryTeam := `
       SELECT team_name
       FROM team
       WHERE team_name = $1
   `
	err := r.tx.QueryRowContext(ctx, queryTeam, name).Scan(&teamID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, team.ErrNotFound
		}
		return nil, err
	}

	queryUsers := `
       SELECT user_id, username, is_active
       FROM "user"
       WHERE team_name = $1
       ORDER BY username
   `

	rows, err := r.tx.QueryContext(ctx, queryUsers, teamID)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}()

	members := make([]team.Member, 0)
	for rows.Next() {
		var m team.Member
		if err := rows.Scan(&m.UserID, &m.Username, &m.IsActive); err != nil {
			return nil, err
		}
		members = append(members, m)
	}

	return &team.Team{
		TeamName: name,
		Members:  members,
	}, nil
}
