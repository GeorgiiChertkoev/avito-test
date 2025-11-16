package repo

import (
	"context"
	"pr-reviewer/internal/domain/pullrequest"
	"pr-reviewer/internal/domain/team"
	"pr-reviewer/internal/domain/user"
)

type Repositories struct {
	TeamRepo team.Repository
	UserRepo user.Repository
	PRRepo   pullrequest.Repository
}

type UnitOfWork interface {
	Do(ctx context.Context, fn func(tx Repositories) error) error
}
