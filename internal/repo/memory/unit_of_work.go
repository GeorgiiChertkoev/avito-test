package memory

import (
	"context"
	"pr-reviewer/internal/domain/repo"
)

type UnitOfWork struct {
	repos repo.Repositories
}

func NewUnitOfWork(repos repo.Repositories) *UnitOfWork {
	return &UnitOfWork{repos: repos}
}

func (u *UnitOfWork) Do(_ context.Context, fn func(repositories repo.Repositories) error) error {
	return fn(u.repos)
}
