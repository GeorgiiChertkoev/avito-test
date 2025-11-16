package memory

import (
	"context"
	"pr-reviewer/internal/domain/pullrequest"
	"sync"
)

type PRRepo struct {
	mu  sync.RWMutex
	prs map[string]*pullrequest.PullRequest
}

func NewPRRepo() *PRRepo {
	return &PRRepo{
		prs: make(map[string]*pullrequest.PullRequest),
	}
}

func (r *PRRepo) CreatePR(_ context.Context, pr *pullrequest.PullRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.prs[pr.ID]; exists {
		return pullrequest.ErrPRExists
	}

	cp := *pr
	r.prs[pr.ID] = &cp

	return nil
}

func (r *PRRepo) GetPR(_ context.Context, id string) (*pullrequest.PullRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	pr, ok := r.prs[id]
	if !ok {
		return nil, pullrequest.ErrNotFound
	}

	cp := *pr

	return &cp, nil
}

func (r *PRRepo) UpdatePR(_ context.Context, pr *pullrequest.PullRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.prs[pr.ID]; !ok {
		return pullrequest.ErrNotFound
	}

	cp := *pr
	r.prs[pr.ID] = &cp

	return nil
}

func (r *PRRepo) GetReviewerPRs(_ context.Context, userID string) ([]pullrequest.PullRequestShort, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]pullrequest.PullRequestShort, 0)

	for _, pr := range r.prs {
		for _, reviewer := range pr.AssignedReviewers {
			if reviewer == userID {
				result = append(result, pullrequest.PullRequestShort{
					ID:       pr.ID,
					Name:     pr.Name,
					AuthorID: pr.AuthorID,
					Status:   pr.Status,
				})

				break
			}
		}
	}

	return result, nil
}
