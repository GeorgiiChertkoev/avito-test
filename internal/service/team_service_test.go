package service_test

import (
	"context"
	"testing"

	"pr-reviewer/internal/domain/repo"
	"pr-reviewer/internal/domain/team"
	"pr-reviewer/internal/repo/memory"
	"pr-reviewer/internal/service"

	"github.com/stretchr/testify/assert"
)

func setupTeamService() (*service.TeamService, repo.Repositories) {
	repos := repo.Repositories{
		UserRepo: memory.NewUserRepo(),
		TeamRepo: memory.NewTeamRepo(),
		PRRepo:   memory.NewPRRepo(),
	}
	uow := memory.NewUnitOfWork(repos)
	teamService := service.NewTeamService(uow)

	return teamService, repos
}

func TestTeamService_CreateTeam(t *testing.T) {
	svc, _ := setupTeamService()

	teamData := &team.Team{
		TeamName: "backend",
		Members: []team.Member{
			{UserID: "u1", Username: "Alice", IsActive: true},
			{UserID: "u2", Username: "Bob", IsActive: true},
		},
	}

	created, err := svc.CreateTeam(context.Background(), teamData)
	assert.NoError(t, err)
	assert.Equal(t, "backend", created.TeamName)
	assert.Len(t, created.Members, 2)
}

func TestTeamService_CreateTeam_AlreadyExists(t *testing.T) {
	svc, _ := setupTeamService()

	teamData := &team.Team{
		TeamName: "backend",
		Members: []team.Member{
			{UserID: "u1", Username: "Alice", IsActive: true},
		},
	}

	// first creation should succeed
	created, err := svc.CreateTeam(context.Background(), teamData)
	assert.NoError(t, err)
	assert.Equal(t, "backend", created.TeamName)

	// second creation should fail
	_, err = svc.CreateTeam(context.Background(), teamData)
	assert.Error(t, err)
	assert.Equal(t, team.ErrTeamExists, err)
}

func TestTeamService_GetTeam(t *testing.T) {
	svc, _ := setupTeamService()

	teamData := &team.Team{
		TeamName: "frontend",
		Members: []team.Member{
			{UserID: "u3", Username: "Charlie", IsActive: true},
		},
	}

	_, err := svc.CreateTeam(context.Background(), teamData)
	assert.NoError(t, err)

	tm, err := svc.GetTeam(context.Background(), "frontend")
	assert.NoError(t, err)
	assert.Equal(t, "frontend", tm.TeamName)
	assert.Len(t, tm.Members, 1)
}

func TestTeamService_GetTeam_NotFound(t *testing.T) {
	svc, _ := setupTeamService()

	tm, err := svc.GetTeam(context.Background(), "nonexistent")
	assert.Nil(t, tm)
	assert.Error(t, err)
}
