package team

import "context"

type Repository interface {
	CreateTeam(ctx context.Context, team *Team) error
	GetTeamByName(ctx context.Context, name string) (*Team, error)
}
