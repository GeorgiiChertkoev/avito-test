package pullrequest

import "context"

type Repository interface {
	CreatePR(ctx context.Context, pr *PullRequest) error
	GetPR(ctx context.Context, id string) (*PullRequest, error)
	UpdatePR(ctx context.Context, pr *PullRequest) error

	GetReviewerPRs(ctx context.Context, userID string) ([]PullRequestShort, error)
}
