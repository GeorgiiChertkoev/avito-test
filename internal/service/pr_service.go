package service

import (
	"context"
	"pr-reviewer/internal/domain/pr"
	"pr-reviewer/internal/domain/repo"
	"pr-reviewer/internal/domain/team"
	"pr-reviewer/internal/domain/user"
	"time"
)

type PRService struct {
	prRepo   pr.Repository
	userRepo user.Repository
	teamRepo team.Repository
	uow      repo.UnitOfWork
}

func NewPRService(prRepo pr.Repository, userRepo user.Repository, teamRepo team.Repository, uow repo.UnitOfWork) *PRService {
	return &PRService{prRepo: prRepo, userRepo: userRepo, teamRepo: teamRepo, uow: uow}
}

type CreatePRInput struct {
	ID       string
	Name     string
	AuthorID string
}

func (s *PRService) CreatePR(ctx context.Context, in CreatePRInput) (*pr.PullRequest, error) {
	var created *pr.PullRequest

	err := s.uow.Do(ctx, func(tx repo.Repositories) error {
		exists, err := tx.PRRepo.PRExists(ctx, in.ID)
		if err != nil {
			return err
		}
		if exists {
			return pr.ErrPRExists
		}

		author, err := tx.UserRepo.GetUserByID(ctx, in.AuthorID)
		if err != nil {
			return pr.ErrPRNotFound
		}

		t, err := tx.TeamRepo.GetTeamByName(ctx, author.TeamName)
		if err != nil {
			return team.ErrTeamNotFound
		}

		candidates := make([]string, 0)
		for _, m := range t.Members {
			if m.IsActive && m.UserID != in.AuthorID {
				candidates = append(candidates, m.UserID)
			}
		}

		assigned := []string{}
		for i := 0; i < len(candidates) && len(assigned) < 2; i++ {
			assigned = append(assigned, candidates[i])
		}

		p := &pr.PullRequest{
			ID:                in.ID,
			Name:              in.Name,
			AuthorID:          in.AuthorID,
			Status:            pr.StatusOpen,
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

func (s *PRService) MergePR(ctx context.Context, id string) (*pr.PullRequest, error) {
	p, err := s.prRepo.GetPR(ctx, id)
	if err != nil {
		return nil, pr.ErrPRNotFound
	}

	if p.Status == pr.StatusMerged {
		return p, nil
	}

	p.Status = pr.StatusMerged
	p.MergedAt = timePtr(time.Now())

	if err := s.prRepo.UpdatePR(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *PRService) ReassignReviewer(ctx context.Context, prID, old string) (*pr.PullRequest, string, error) {
	//var replacedBy string

	p, err := s.prRepo.GetPR(ctx, prID)
	if err != nil {
		return nil, "", pr.ErrPRNotFound
	}

	if p.Status == pr.StatusMerged {
		return nil, "", pr.ErrPRMerged
	}

	assigned := p.AssignedReviewers
	idx := -1
	for i, r := range assigned {
		if r == old {
			idx = i
			break
		}
	}
	if idx == -1 {
		return nil, "", pr.ErrNotAssigned
	}

	author, err := s.userRepo.GetUserByID(ctx, p.AuthorID)
	if err != nil {
		return nil, "", pr.ErrPRNotFound
	}

	t, err := s.teamRepo.GetTeamByName(ctx, author.TeamName)
	if err != nil {
		return nil, "", team.ErrTeamNotFound
	}

	candidates := []string{}
	for _, m := range t.Members {
		if m.IsActive && m.UserID != old && m.UserID != author.UserID {
			candidates = append(candidates, m.UserID)
		}
	}

	if len(candidates) == 0 {
		return nil, "", pr.ErrNoCandidate
	}

	replacement := candidates[0]
	assigned[idx] = replacement
	p.AssignedReviewers = assigned

	if err := s.prRepo.UpdatePR(ctx, p); err != nil {
		return nil, "", err
	}

	return p, replacement, nil
}

func timePtr(t time.Time) *time.Time { return &t }
