package pr

import "context"

type Repository interface {
	CreatePR(ctx context.Context, pr *PullRequest) error
	GetPR(ctx context.Context, id string) (*PullRequest, error)
	UpdatePR(ctx context.Context, pr *PullRequest) error

	PRExists(ctx context.Context, id string) (bool, error)
}
