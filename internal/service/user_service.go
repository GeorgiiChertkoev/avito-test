package service

import (
	"context"
	"pr-reviewer/internal/domain/pr"
	"pr-reviewer/internal/domain/user"
)

type UserService struct {
	userRepo user.Repository
	prRepo   pr.Repository
}

func NewUserService(userRepo user.Repository, prRepo pr.Repository) *UserService {
	return &UserService{userRepo: userRepo, prRepo: prRepo}
}

func (s *UserService) SetIsActive(ctx context.Context, userID string, active bool) (*user.User, error) {
	u, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, user.ErrUserNotFound
	}

	u.IsActive = active
	updated, err := s.userRepo.SetIsActive(ctx, userID, active)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

type UserReviews struct {
	UserID       string                `json:"user_id"`
	PullRequests []pr.PullRequestShort `json:"pull_requests"`
}

func (s *UserService) GetReviewPRs(ctx context.Context, userID string) (*UserReviews, error) {
	_, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, user.ErrUserNotFound
	}

	prs, err := s.userRepo.GetReviewerPRs(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &UserReviews{
		UserID:       userID,
		PullRequests: prs,
	}, nil
}
