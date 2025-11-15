package user

import (
	"context"
	"pr-reviewer/internal/domain/pr"
)

type Repository interface {
	UpsertUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id string) (*User, error)
	SetIsActive(ctx context.Context, id string, active bool) (*User, error)

	GetUsersByTeam(ctx context.Context, teamName string) ([]User, error)

	// Domain-oriented query:
	// PRs where this user is a reviewer (PRShort models live in pr domain)
	GetReviewerPRs(ctx context.Context, userID string) ([]pr.PullRequestShort, error)
}
