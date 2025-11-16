package service

import (
	"context"
	"pr-reviewer/internal/domain/pullrequest"
	"pr-reviewer/internal/domain/repo"
	"pr-reviewer/internal/domain/user"
)

type UserService struct {
	uow repo.UnitOfWork
}

func NewUserService(uow repo.UnitOfWork) *UserService {
	return &UserService{uow: uow}
}

func (s *UserService) SetIsActive(ctx context.Context, userID string, active bool) (*user.User, error) {
	var usr *user.User

	err := s.uow.Do(ctx, func(tx repo.Repositories) error {
		var err error
		usr, err = tx.UserRepo.GetUserByID(ctx, userID)
		if err != nil {
			return err
		}

		usr.IsActive = active
		err = tx.UserRepo.UpsertUser(ctx, usr)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return usr, nil
}

type UserReviews struct {
	UserID       string                         `json:"user_id"`
	PullRequests []pullrequest.PullRequestShort `json:"pull_requests"`
}

func (s *UserService) GetReviewPRs(ctx context.Context, userID string) (*UserReviews, error) {
	var prs []pullrequest.PullRequestShort

	err := s.uow.Do(ctx, func(tx repo.Repositories) error {
		var err error
		prs, err = tx.PRRepo.GetReviewerPRs(ctx, userID)

		return err
	})
	if err != nil {
		return nil, err
	}

	return &UserReviews{
		UserID:       userID,
		PullRequests: prs,
	}, nil
}
