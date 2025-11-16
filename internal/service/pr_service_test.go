package service_test

import (
	"context"
	"testing"
	"time"

	"pr-reviewer/internal/domain/pullrequest"
	"pr-reviewer/internal/domain/repo"
	"pr-reviewer/internal/domain/team"
	"pr-reviewer/internal/domain/user"
	"pr-reviewer/internal/repo/memory"
	"pr-reviewer/internal/service"

	"github.com/stretchr/testify/require"
)

func setupPRService() (*service.PRService, repo.Repositories) {
	repos := repo.Repositories{
		UserRepo: memory.NewUserRepo(),
		TeamRepo: memory.NewTeamRepo(),
		PRRepo:   memory.NewPRRepo(),
	}
	uow := memory.NewUnitOfWork(repos)
	prService := service.NewPRService(uow)

	return prService, repos
}

func populateTeam(repos repo.Repositories, teamName string, members []user.User) {
	_ = repos.TeamRepo.CreateTeam(context.Background(), &team.Team{
		TeamName: teamName,
		Members: func() []team.Member {
			tm := make([]team.Member, len(members))
			for i, m := range members {
				tm[i] = team.Member{
					UserID:   m.UserID,
					Username: m.Username,
					IsActive: m.IsActive,
				}
			}

			return tm
		}(),
	})

	for _, m := range members {
		_ = repos.UserRepo.UpsertUser(context.Background(), &m)
	}
}

func TestPRService_CreatePR(t *testing.T) {
	prService, repos := setupPRService()

	populateTeam(repos, "backend", []user.User{
		{UserID: "u1", Username: "Alice", TeamName: "backend", IsActive: true},
		{UserID: "u2", Username: "Bob", TeamName: "backend", IsActive: true},
		{UserID: "u3", Username: "Charlie", TeamName: "backend", IsActive: false},
	})

	input := service.CreatePRInput{
		ID:       "pr-1",
		Name:     "Add search",
		AuthorID: "u1",
	}

	pr, err := prService.CreatePR(context.Background(), input)
	require.NoError(t, err)
	require.Equal(t, "pr-1", pr.ID)
	require.Equal(t, "Add search", pr.Name)
	require.Equal(t, pullrequest.StatusOpen, pr.Status)
	require.Len(t, pr.AssignedReviewers, 1) // only u2 is active
	require.Contains(t, pr.AssignedReviewers, "u2")
}

func TestPRService_MergePR(t *testing.T) {
	prService, repos := setupPRService()

	_ = repos.PRRepo.CreatePR(context.Background(), &pullrequest.PullRequest{
		ID:                "pr-1",
		Name:              "Fix bug",
		AuthorID:          "u1",
		Status:            pullrequest.StatusOpen,
		AssignedReviewers: []string{"u2"},
		CreatedAt:         timePtr(time.Now()),
	})

	pr, err := prService.MergePR(context.Background(), "pr-1")
	require.NoError(t, err)
	require.Equal(t, pullrequest.StatusMerged, pr.Status)
	require.NotNil(t, pr.MergedAt)

	// merging again is idempotent
	pr2, err := prService.MergePR(context.Background(), "pr-1")
	require.NoError(t, err)
	require.Equal(t, pullrequest.StatusMerged, pr2.Status)
}

func TestPRService_ReassignReviewer(t *testing.T) {
	prService, repos := setupPRService()

	populateTeam(repos, "backend", []user.User{
		{UserID: "u1", Username: "Alice", TeamName: "backend", IsActive: true},
		{UserID: "u2", Username: "Bob", TeamName: "backend", IsActive: true},
		{UserID: "u3", Username: "Charlie", TeamName: "backend", IsActive: true},
	})

	_ = repos.PRRepo.CreatePR(context.Background(), &pullrequest.PullRequest{
		ID:                "pr-1",
		Name:              "Feature X",
		AuthorID:          "u1",
		Status:            pullrequest.StatusOpen,
		AssignedReviewers: []string{"u2"},
		CreatedAt:         timePtr(time.Now()),
	})

	pr, replacedBy, err := prService.ReassignReviewer(context.Background(), "pr-1", "u2")
	require.NoError(t, err)
	require.Equal(t, 1, len(pr.AssignedReviewers))
	require.Equal(t, replacedBy, pr.AssignedReviewers[0])
	require.NotEqual(t, "u2", replacedBy)
}

func TestPRService_ReassignReviewer_Errors(t *testing.T) {
	prService, repos := setupPRService()

	populateTeam(repos, "backend", []user.User{
		{UserID: "u1", Username: "Alice", TeamName: "backend", IsActive: true},
		{UserID: "u2", Username: "Bob", TeamName: "backend", IsActive: true},
	})

	// PR does not exist
	_, _, err := prService.ReassignReviewer(context.Background(), "pr-unknown", "u2")
	require.Error(t, err)

	// PR merged
	_ = repos.PRRepo.CreatePR(context.Background(), &pullrequest.PullRequest{
		ID:                "pr-1",
		Name:              "Feature X",
		AuthorID:          "u1",
		Status:            pullrequest.StatusMerged,
		AssignedReviewers: []string{"u2"},
	})
	_, _, err = prService.ReassignReviewer(context.Background(), "pr-1", "u2")
	require.Equal(t, pullrequest.ErrPRMerged, err)

	// reviewer not assigned
	_ = repos.PRRepo.CreatePR(context.Background(), &pullrequest.PullRequest{
		ID:                "pr-2",
		Name:              "Feature Y",
		AuthorID:          "u1",
		Status:            pullrequest.StatusOpen,
		AssignedReviewers: []string{"u2"},
	})
	_, _, err = prService.ReassignReviewer(context.Background(), "pr-2", "u3")
	require.Equal(t, pullrequest.ErrNotAssigned, err)
}

func timePtr(t time.Time) *time.Time { return &t }
