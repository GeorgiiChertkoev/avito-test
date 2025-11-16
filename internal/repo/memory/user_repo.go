package memory

import (
	"context"
	"pr-reviewer/internal/domain/user"
	"sync"
)

type UserRepo struct {
	mu    sync.RWMutex
	users map[string]*user.User
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		users: make(map[string]*user.User),
	}
}

func (u *UserRepo) UpsertUser(_ context.Context, usr *user.User) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	clone := *usr
	u.users[usr.UserID] = &clone

	return nil
}

func (u *UserRepo) GetUserByID(_ context.Context, id string) (*user.User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	usr, ok := u.users[id]
	if !ok {
		return nil, user.ErrNotFound
	}

	clone := *usr

	return &clone, nil
}

func (u *UserRepo) SetIsActive(_ context.Context, id string, active bool) (*user.User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	usr, ok := u.users[id]
	if !ok {
		return nil, user.ErrNotFound
	}

	usr.IsActive = active
	clone := *usr

	return &clone, nil
}
