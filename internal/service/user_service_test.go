package service_test

import (
	"context"
	"testing"

	"pr-reviewer/internal/domain/pullrequest"
	"pr-reviewer/internal/domain/repo"
	"pr-reviewer/internal/domain/user"
	"pr-reviewer/internal/repo/memory"
	"pr-reviewer/internal/service"

	"github.com/stretchr/testify/require"
)

func setupUserService() (*service.UserService, repo.Repositories) {
	repos := repo.Repositories{
		UserRepo: memory.NewUserRepo(),
		TeamRepo: memory.NewTeamRepo(),
		PRRepo:   memory.NewPRRepo(),
	}

	uow := memory.NewUnitOfWork(repos)
	userService := service.NewUserService(uow)

	return userService, repos
}

func TestUserService_SetIsActive(t *testing.T) {
	svc, repos := setupUserService()

	// Prepopulate user
	u := &user.User{
		UserID:   "u1",
		Username: "Alice",
		TeamName: "backend",
		IsActive: true,
	}
	_ = repos.UserRepo.UpsertUser(context.Background(), u)

	// Deactivate user
	updated, err := svc.SetIsActive(context.Background(), "u1", false)
	require.NoError(t, err)
	require.False(t, updated.IsActive)

	// Reactivate user
	updated, err = svc.SetIsActive(context.Background(), "u1", true)
	require.NoError(t, err)
	require.True(t, updated.IsActive)

	// Non-existent user
	updated, err = svc.SetIsActive(context.Background(), "u999", true)
	require.Error(t, err)
	require.Nil(t, updated)
}

func TestUserService_GetReviewPRs(t *testing.T) {
	svc, repos := setupUserService()

	ctx := context.Background()

	// Prepopulate PRs
	pr1 := &pullrequest.PullRequest{
		ID:                "pr-1",
		Name:              "Add search",
		AuthorID:          "u1",
		Status:            pullrequest.StatusOpen,
		AssignedReviewers: []string{"u2", "u3"},
	}

	pr2 := &pullrequest.PullRequest{
		ID:                "pr-2",
		Name:              "Fix bug",
		AuthorID:          "u2",
		Status:            pullrequest.StatusOpen,
		AssignedReviewers: []string{"u2"},
	}

	_ = repos.PRRepo.CreatePR(ctx, pr1)
	_ = repos.PRRepo.CreatePR(ctx, pr2)

	reviews, err := svc.GetReviewPRs(ctx, "u2")
	require.NoError(t, err)
	require.Len(t, reviews.PullRequests, 2)

	reviews, err = svc.GetReviewPRs(ctx, "u3")
	require.NoError(t, err)
	require.Len(t, reviews.PullRequests, 1)
	require.Equal(t, "pr-1", reviews.PullRequests[0].ID)

	// Test user with no PRs
	reviews, err = svc.GetReviewPRs(ctx, "u999")
	require.NoError(t, err)
	require.Len(t, reviews.PullRequests, 0)
}
