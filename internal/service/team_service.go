package service

import (
	"context"
	"pr-reviewer/internal/domain/repo"
	"pr-reviewer/internal/domain/team"
	"pr-reviewer/internal/domain/user"
)

type TeamService struct {
	teamRepo team.Repository
	userRepo user.Repository
	uow      repo.UnitOfWork
}

func NewTeamService(teamRepo team.Repository, userRepo user.Repository, ouw repo.UnitOfWork) *TeamService {
	return &TeamService{teamRepo: teamRepo, userRepo: userRepo, uow: ouw}
}

func (s *TeamService) AddTeam(ctx context.Context, t *team.Team) (*team.Team, error) {
	var created *team.Team

	err := s.uow.Do(ctx, func(tx repo.Repositories) error {
		existing, err := tx.TeamRepo.GetTeamByName(ctx, t.TeamName)
		if err != nil {
			return err
		}
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
	t, err := s.teamRepo.GetTeamByName(ctx, name)
	if err != nil {
		return nil, team.ErrTeamNotFound
	}
	return t, nil
}
