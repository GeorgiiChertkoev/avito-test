package memory

import (
	"context"
	"log"
	"pr-reviewer/internal/domain/team"
	"sync"
)

type TeamRepo struct {
	mu    sync.RWMutex
	teams map[string]*team.Team
}

func NewTeamRepo() *TeamRepo {
	return &TeamRepo{
		teams: make(map[string]*team.Team),
	}
}

func (r *TeamRepo) CreateTeam(_ context.Context, t *team.Team) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.teams[t.TeamName]; exists {
		log.Printf("TeamRepo.CreateTeam: Team already exists: %s", t.TeamName)
		return team.ErrTeamExists
	}
	r.teams[t.TeamName] = t

	return nil
}

func (r *TeamRepo) GetTeamByName(_ context.Context, name string) (*team.Team, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, exists := r.teams[name]
	if !exists {
		return nil, team.ErrNotFound
	}

	return t, nil
}
