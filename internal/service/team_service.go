package service

import (
	"context"
	"pr-reviewer/internal/domain/repo"
	"pr-reviewer/internal/domain/team"
	"pr-reviewer/internal/domain/user"
)

type TeamService struct {
	uow repo.UnitOfWork
}

func NewTeamService(ouw repo.UnitOfWork) *TeamService {
	return &TeamService{uow: ouw}
}

func (s *TeamService) CreateTeam(ctx context.Context, t *team.Team) (*team.Team, error) {
	var created *team.Team

	err := s.uow.Do(ctx, func(tx repo.Repositories) error {
		existing, _ := tx.TeamRepo.GetTeamByName(ctx, t.TeamName)
		if existing != nil {
			return team.ErrTeamExists
		}

		if err := tx.TeamRepo.CreateTeam(ctx, t); err != nil {
			return err
		}

		for i := range t.Members {
			m := &t.Members[i]
			u := &user.User{
				UserID:   m.UserID,
				Username: m.Username,
				TeamName: t.TeamName,
				IsActive: m.IsActive,
			}
			if err := tx.UserRepo.UpsertUser(ctx, u); err != nil {
				return err
			}
		}
		created = t

		return nil
	})

	return created, err
}

func (s *TeamService) GetTeam(ctx context.Context, name string) (*team.Team, error) {
	var t *team.Team
	err := s.uow.Do(ctx, func(tx repo.Repositories) error {
		var err error
		t, err = tx.TeamRepo.GetTeamByName(ctx, name)

		return err
	})
	if err != nil {
		return nil, err
	}

	return t, nil
}
