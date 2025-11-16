package user

import (
	"context"
)

type Repository interface {
	UpsertUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id string) (*User, error)
	SetIsActive(ctx context.Context, id string, active bool) (*User, error)
}
