package repo

import (
	"context"
	"pr-reviewer/internal/domain/pr"
	"pr-reviewer/internal/domain/team"
	"pr-reviewer/internal/domain/user"
)

type Repositories struct {
	TeamRepo team.Repository
	UserRepo user.Repository
	PRRepo   pr.Repository
}

type UnitOfWork interface {
	Do(ctx context.Context, fn func(tx Repositories) error) error
}
