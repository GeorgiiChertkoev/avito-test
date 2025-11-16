package service

import (
	"context"
	"pr-reviewer/internal/domain/pullrequest"
	"pr-reviewer/internal/domain/repo"
	"pr-reviewer/internal/domain/team"
	"slices"
	"time"
)

type PRService struct {
	uow repo.UnitOfWork
}

func NewPRService(uow repo.UnitOfWork) *PRService {
	return &PRService{uow: uow}
}

type CreatePRInput struct {
	ID       string
	Name     string
	AuthorID string
}

func (s *PRService) CreatePR(ctx context.Context, in CreatePRInput) (*pullrequest.PullRequest, error) {
	var created *pullrequest.PullRequest

	err := s.uow.Do(ctx, func(tx repo.Repositories) error {
		t, err := getTeamByUserId(ctx, tx, in.AuthorID)
		if err != nil {
			return err
		}

		assigned := pickReviewers(t, 2, []string{in.AuthorID})

		p := &pullrequest.PullRequest{
			ID:                in.ID,
			Name:              in.Name,
			AuthorID:          in.AuthorID,
			Status:            pullrequest.StatusOpen,
			AssignedReviewers: assigned,
			CreatedAt:         timePtr(time.Now()),
		}

		if err := tx.PRRepo.CreatePR(ctx, p); err != nil {
			return err
		}
		created = p

		return nil
	})

	return created, err
}

func (s *PRService) MergePR(ctx context.Context, id string) (*pullrequest.PullRequest, error) {
	var pr *pullrequest.PullRequest
	err := s.uow.Do(context.Background(), func(tx repo.Repositories) error {
		var err error
		pr, err = tx.PRRepo.GetPR(ctx, id)
		if err != nil {
			return err
		}

		if pr.Status == pullrequest.StatusMerged {
			return nil
		}

		pr.Status = pullrequest.StatusMerged
		pr.MergedAt = timePtr(time.Now())

		err = tx.PRRepo.UpdatePR(ctx, pr)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *PRService) ReassignReviewer(ctx context.Context, prID, oldReviewer string) (*pullrequest.PullRequest, string, error) {
	var pr *pullrequest.PullRequest
	var replacedBy string
	err := s.uow.Do(context.Background(), func(tx repo.Repositories) error {
		var err error
		pr, err = tx.PRRepo.GetPR(ctx, prID)
		if err != nil {
			return err
		}
		if pr.Status == pullrequest.StatusMerged {
			return pullrequest.ErrPRMerged
		}

		if !slices.Contains(pr.AssignedReviewers, oldReviewer) {
			return pullrequest.ErrNotAssigned
		}

		t, err := getTeamByUserId(ctx, tx, oldReviewer)
		if err != nil {
			return err
		}

		reassigned := pickReviewers(t, 1, pr.AssignedReviewers)
		if len(reassigned) == 0 {
			return pullrequest.ErrNoCandidate
		}
		replacedBy = reassigned[0]

		pr.AssignedReviewers = slices.DeleteFunc(pr.AssignedReviewers, func(id string) bool {
			return id == oldReviewer
		})
		pr.AssignedReviewers = append(pr.AssignedReviewers, reassigned...)

		err = tx.PRRepo.UpdatePR(ctx, pr)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, "", err
	}

	return pr, replacedBy, nil
}

func timePtr(t time.Time) *time.Time { return &t }

func getTeamByUserId(ctx context.Context, tx repo.Repositories, userId string) (*team.Team, error) {
	author, err := tx.UserRepo.GetUserByID(ctx, userId)
	if err != nil {
		return nil, pullrequest.ErrNotFound
	}

	t, err := tx.TeamRepo.GetTeamByName(ctx, author.TeamName)
	if err != nil {
		return nil, pullrequest.ErrNotFound
	}

	return t, nil
}

func pickReviewers(t *team.Team, reviewerCount int, excluded []string) []string {
	candidates := make([]string, 0)
	for _, m := range t.Members {
		if m.IsActive && !slices.Contains(excluded, m.UserID) {
			candidates = append(candidates, m.UserID)
		}
	}

	assigned := make([]string, 0)
	for i := 0; i < len(candidates) && len(assigned) < reviewerCount; i++ {
		assigned = append(assigned, candidates[i])
	}

	return assigned
}
